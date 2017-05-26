package utils

import (
	"fmt"
)

const (
	// FmtPrefix to be included in all log traces
	FmtPrefix = "== wifi-connect"
)

// Errorf returns a new instance implementing error interface taken a formatted string and
// params as input
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf("%v: %v", FmtPrefix, format), a...)
}

// PkgErrorf returns a new instance implementing error interface taken a formatted string and
// params as input
func PkgErrorf(pkg string, format string, a ...interface{}) error {
	return fmt.Errorf(fmt.Sprintf("%v/%v: %v", FmtPrefix, pkg, format), a...)
}

// Sprintf returns a formatted string with project prefix
func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%v: %v", FmtPrefix, format), a...)
}

// PkgSprintf returns a formatted string with project prefix
func PkgSprintf(pkg string, format string, a ...interface{}) string {
	return fmt.Sprintf(fmt.Sprintf("%v/%v: %v", FmtPrefix, pkg, format), a...)
}
