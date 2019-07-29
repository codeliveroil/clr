package test

import (
	"reflect"
	"testing"
)

// Compare compares got and want and logs an error if they
// are not equal. The line number of the callee will be
// present in the error log.
func Compare(t *testing.T, got interface{}, want interface{}) {

	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got: %v, want: %v", got, want)
	}

}

// Err logs an error if err is not nil. The line number of
// the callee will be present in the error log.
func Err(t *testing.T, err error) {

	t.Helper()

	if err != nil {
		t.Error(err)
	}

}

// Fatal logs an error if err is not nil and fatals. The line
// number of the callee will be present in the error log.
func Fatal(t *testing.T, err error) {

	t.Helper()

	if err != nil {
		t.Fatal(err)
	}

}
