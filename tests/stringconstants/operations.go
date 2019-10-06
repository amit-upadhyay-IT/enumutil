package stringconstants

import (
	"enumutils/enumutil"
	"fmt"
	"log"
	"runtime"
)

type Operation string

const (
	ALIAS   Operation = "-a"    // normal alias
	ALIAS_P Operation = "-ap"   // parameterized alias
	ALIAS_D Operation = "-d"    // to delete both type of alias(i.e. normal and parameterized)
	ALIAS_U Operation = "-u"    // to update the keys or values for an alias
	STATUS  Operation = "-s"    // to see all the alias added
	RETRY   Operation = "retry" // incase wrong input is entered by user
	HELP    Operation = "-h"    // show possible commands with details

	PERFORM Operation = "p" // when user will enter command to perform task,
	SOMETHINGELSE
	AMIT1 string = "amit1"
	WHATEVER
	// the parser should be smart enough to figure out if the user is trying
	// to perform a valid operation, else parser should return RETRY flag.
)

func GetEnums() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Failed to recover the file name")
	}
	enum := enumutil.Enum()
	enum.FetchEnums(file)
	fmt.Println(enum)
}
