package utils

import (
	"testing"
)

func TestErrorFormat(t *testing.T) {

	msg := "Whatever message"
	e := Errorf(msg)
	expected := FmtPrefix + ": " + msg

	if e.Error() != expected {
		t.Errorf("Error is not well formatted, expected %v but got %v", expected, e.Error())
	}
}
func TestErrorfFormat(t *testing.T) {

	format := "Whatever %v"
	detail := "message"
	e := Errorf(format, detail)

	expected := FmtPrefix + ": Whatever message"

	if e.Error() != expected {
		t.Errorf("Error is not well formatted, expected %s but got %s", expected, e.Error())
	}
}

func TestSprintFormat(t *testing.T) {

	msg := "Whatever message"
	e := Sprintf(msg)
	expected := FmtPrefix + ": " + msg

	if e != expected {
		t.Errorf("String is not well formatted, expected %s but got %s", expected, e)
	}
}

func TestSprintfFormat(t *testing.T) {

	format := "Whatever %v"
	detail := "message"
	e := Sprintf(format, detail)

	expected := FmtPrefix + ": Whatever message"

	if e != expected {
		t.Errorf("String is not well formatted, expected %s but got %s", expected, e)
	}
}
