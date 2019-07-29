package errors

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

// Exit exits with an error code 1 if err is not nil. It
// prints a friendly message and the line number where
// the error occurs.
func Exit(err error, format string, args ...interface{}) {
	if err == nil {
		return
	}

	var loc string
	if _, file, line, ok := runtime.Caller(1); ok {
		loc = fmt.Sprintf("[%s:%d] ", path.Base(file), line)
	}

	args = append(args, err)
	fmt.Printf(loc+format+": %v\n", args...)
	os.Exit(1)
}
