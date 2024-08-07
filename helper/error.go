package helper

import (
	"fmt"
	"runtime"
	"strings"
)

func Error(message error) error {
	pc, file, line, _ := runtime.Caller(1)

	// File path
	splitFiles := strings.Split(file, "/internal/")
	shortFile := splitFiles[0]
	if len(splitFiles) > 0 {
		shortFile = file[len(shortFile)+10:]
	}

	// Function name
	fullFnName := runtime.FuncForPC(pc).Name()
	splitFnName := strings.Split(fullFnName, "/internal/")
	fnName := splitFnName[0]
	if len(splitFnName) > 0 {
		fnName = fullFnName[len(fnName)+10:]
	}

	return fmt.Errorf("%s \n...%s:%d (%s)", message, shortFile, line, fnName)
}
