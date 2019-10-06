package enumutil

/*
Generally in Go, an enum is represented using constant variables of some primitive type, (they can't be non-primitive,
however their name can be aliased, eg: `type CustomType int`)
This is a utility package for enums, currently, it is maintaining two map and other variables:
1) enumMapsStr := map[string]map[string]string : key is the enum name, value is again a map of string to string, where
	value's key is constant variable name and value's value to constant variable value.
2) enumMapsInt := map[string]map[string]int : key is the enum name, value is a map of string to int, where values'key
	is constant variable name and value's value is constant variable value.
As you can see the above maps, they are different for storing different values types. Thus, we also need a mechanism to
put the right enum to the right map. So, for that I am maintaining one more map.
3) typeAliasMap := map[string]string : this map will store the type alias to its type value, eg: `type Operation string`
	typeAliasMap["Operation"] = "int".
Validations are very crucial part here, validating the typeAliasMap to always have a primitive type if its being used as
a constant variable, etc.
 */

import (
	"enumutils/utils/io"
	"log"
	"strings"
	"sync"
)

type enum struct {
	enumMapStrStore map[string]map[string]string
	enumMapIntStore map[string]map[string]int
	typeAliasMap map[string]string
}

var enumInstance *enum = nil
var once sync.Once

func Enum() *enum {
	once.Do(func() {
		enumInstance = &enum{
			enumMapStrStore: make(map[string]map[string]string),
			enumMapIntStore: make(map[string]map[string]int),
			typeAliasMap:    make(map[string]string),
		}
	})
	return enumInstance
}

func (e *enum) GetIntegerEnums() map[string]map[string]int {
	return e.enumMapIntStore
}

func (e * enum) GetStringEnums() map[string]map[string]string {
	return e.enumMapStrStore
}

func (e *enum) getTypeAliasMap() map[string]string {
	return e.typeAliasMap
}

// storing the last updated enum contents
type prevEnumInfos struct {
	varType  string // can be string, int, or etc
	varValue string // value of const var
	varEnum  string // category to which const var was declared
}

func (e *enum) FetchEnums(filename string) {

	// read contents from file
	lines, err := io.ReadFile(filename, true)
	if err != nil {
		log.Fatal(err)
	}

	// initialize the typeAliasMap, because it's going to be used in the algorithm below
	e.typeAliasMap = getTypeAliasMap(lines)

	// iterate over file contents
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "const(") || strings.HasPrefix(lines[i], "const (") {
			prevInfos := &prevEnumInfos{} // this previous information is one for each const block, coz other const blocks needs fresh info. So not keeping it as a member of Enum class
			i = e.fetchEnumsFromBlock(lines, i, prevInfos)
		}
	}
}

// reads a const block
func (e *enum) fetchEnumsFromBlock(lines []string, begIndex int, prevInfos *prevEnumInfos) int { // maybe later we can return error also(if found)
	var parenthesisCount int
	parenthesisCount = 1
	i := begIndex + 1
	// iterating until const block is closed, NOTE: the comments inside const block should have balanced parenthesis
	for parenthesisCount != 0 {
		line := lines[i]
		words := strings.Fields(line)
		// check for parenthesis
		if strings.Contains(line, "(") {
			parenthesisCount++
		}
		if strings.Contains(line, ")") {
			parenthesisCount--
		}
		// continue adding enum values
		if len(words) > 0 {
			e.fetchEnumsFromLine(line, prevInfos)
		}
		i++
	}
	return i
}

// continue adding enum values
// rules:
// 1) Line should not be a comment or should not start with )
// 2) If words has length more than 1, then second word should either match the enumName or it should be a comment, if its neither of these then we assume its a declaration of some other constant.
// 3) If its declaration of some other constant, mark a flag representing the same
// 4) if words length is 1
// 		4.1) flag declared in point 3 is false, then simply add it to enum values
// 		4.2) flag declared in point 3 is true, then create a entry in the map with the second word as the key and add the first word in the line as the value for that key
func (e *enum) fetchEnumsFromLine(line string, prevInfos *prevEnumInfos) {

	// remove comments from the line
	if strings.Contains(line, "//") {
		parts := strings.Split(line, "//")
		line = parts[0]
	}

	words := strings.Fields(line)
	if len(words) > 0 && !strings.HasPrefix(words[0], ")") { // Ignoring comments or closing const block, de-moran's law: !A || !B ~ !(A && B)
		if len(words) == 1 {
			// cases:
			// 1) type can be string:
			// 		1.1) simply set same type and value to the current const variable
			// 2) type can be int
			// I am not supporting other types now
			if prevInfos.varType == "string" {
				if prevInfos.varEnum != "" {
					constVarsMap, ok := e.enumMapStrStore[prevInfos.varEnum]
					if !ok {
						log.Fatal("Seems wrong, not possible to have the enum type in enumStringMapStore")
					}
					constVarsMap[words[0]] = prevInfos.varValue;
					// set values in prevInfos, in case len(words) is 1, the type and enum name will be same so not setting them
					prevInfos.varValue = constVarsMap[words[0]]
				} else {
					log.Fatal("this seems wrong, as no const var can be declared without a type")
				}
			} else if prevInfos.varType == "int" {
				// its value can be either generated from an expression(we will need to maintain an abstract syntax tree) or it can be default(int)
			}
		} else if len(words) == 2 {
			// eg: QUEUED State="queued", has len(words) as 2
			// or PROCESSED int=1, also have len(words) as 2, but I need to ignore it
			if strings.Contains(words[1], "=") {
				splits := strings.Split(words[1], "=")
				enumName := splits[0]
				varValue := splits[1]
				e.addEnumDetailsToStrMap(enumName, varValue, words, prevInfos)
			} else {
				log.Fatal("Seems wrong, as this would be a syntax error")
			}
		} else if len(words) == 3 {
			// eg: PHI = 1.618 has 3 words, so we will ignore it as its type is float and someone using enum should name it using type alias
			// or PROCESSING State= "processing" or PROCESSING State ="processing"
			enumName := ""
			if strings.Contains(words[1], "=") {
				enumName = strings.Trim(words[1], "=")
			} else if strings.Contains(words[2], "=") {
				enumName = strings.Trim(words[2], "=")
			} else {
				// ignoring this case
				return
			}
			varValue := words[2]
			e.addEnumDetailsToStrMap(enumName, varValue, words, prevInfos)
		} else if len(words) == 4 {
			// words[0] should be const var name
			// words[1] should be var type
			// words[2] should be =
			// words[3] should be value
			enumName := words[1]
			varValue := words[3]
			e.addEnumDetailsToStrMap(enumName, varValue, words, prevInfos)
		} else if len(words) >= 5 {
			// this case can only be present when const vars has number type(int or float), because then it would be
			// an expression which needs to be evaluated with abstract syntax tree,
			// eg: KB ByteSize = 1 << (10 * iota), here len(ords) = 8
			// TODO: implement it
		}
	}
}

func (e *enum) addEnumDetailsToStrMap(enumName, varValue string, words []string, prevInfos *prevEnumInfos) {

	// remove " from varValue
	varValue = strings.Trim(varValue, "\"")
	// ignore if enumName is not part of e.typeAliasMap
	if varType, ok := e.typeAliasMap[enumName]; ok {
		if varType == "string" {
			if prevInfos.varEnum != "" {
				constVarsMap, has := e.enumMapStrStore[prevInfos.varEnum]
				if has {
					constVarsMap[words[0]] = varValue
				} else {
					temp := make(map[string]string)
					temp[words[0]] = varValue
					e.enumMapStrStore[enumName] = temp
				}
			} else {
				temp := make(map[string]string)
				temp[words[0]] = varValue
				e.enumMapStrStore[enumName] = temp
			}
			// set up the values in prevInfos
			prevInfos.varType = "string" //we can remove its usage coz its same as e.typeAliasMap[prevInfos.varEnum]
			prevInfos.varValue = varValue
			prevInfos.varEnum = enumName
		} else if varType == "int" { //TODO: separate it and make different method for integer
			// TODO
		}
	}
}

func getTypeAliasMap(lines []string) map[string]string {
	m := make(map[string]string)
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		words := strings.Fields(line)
		if len(words) >= 3 && words[0] == "type" && isPrimitive(words[2]) {
			m[words[1]] = words[2]
		}
	}
	return m
}

func isPrimitive(word string) bool {
	primitives := []string {"bool",
		"string",
		"int","int8" ,"int16" ,"int32" ,"int64",
		"uint","uint8","uint16","uint32","uint64","uintptr",
		"byte",
		"rune",
		"float32","float64",
		"complex64","complex128"}
	for _, item := range primitives {
		if item == word {
			return true
		}
	}
	return false
}
