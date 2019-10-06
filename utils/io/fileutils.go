package io

import (
	"bufio"
	"os"
	"os/exec"
	"strings"
)

/*
* File Input/Output utility functions
* - CreateFile(string) error
* - ReadFile(string) []string, error
* - AppendToFile(string, string, string) error
 */

func IsFilePresent(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// creates a file with the path specified
// TODO: what if path is not relative? will it create file in current dir?
func CreateFile(path string) error {
	// detect if file exists
	//var _, err = os.Stat(path)
	//if err != nil {
	//    return err
	//}

	file, err := os.Create(path)
	if err != nil {
		// create file using mkdir and touch command
		err1 := createFileUsingMkdir(path)
		if err1 != nil {
			return err1
		}
		return err
	}
	defer file.Close()
	return nil
}

// reads file and returns slice of string
// TODO: what would happen if the file is very long?
// would the slice of string be able to read all the content in one go?
func ReadFile(name string, removeWhiteSpaces bool) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		file.Close()
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	fileLines := []string{}
	line := ""
	for scanner.Scan() {
		line = scanner.Text()
		if removeWhiteSpaces {
			line = removeExtraSpaces(line)
		}
		fileLines = append(fileLines, line)
	}
	return fileLines, nil
}

func AppendToFile(filename, key, value string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		file.Close()
		return err
	}
	defer file.Close()

	if _, err := file.Write([]byte(key + "," + value)); err != nil {
		return err
	}
	return nil
}

func createFileUsingMkdir(filepath string) error {
	cmd := exec.Command("mkdir", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func removeExtraSpaces(token string) string {
	token = strings.TrimSpace(token)
	token = strings.TrimRight(token, "\n")
	return token
}
