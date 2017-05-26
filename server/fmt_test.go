package server

import (
	"testing"

	"github.com/CanonicalLtd/UCWifiConnect/utils"
)

const (
	pkg = "server"
)

func TestErrorFormat(t *testing.T) {

	msg := "Whatever message"
	e := Errorf(msg)
	expected := utils.FmtPrefix + "/" + pkg + ": " + msg

	if e.Error() != expected {
		t.Errorf("Error is not well formatted, expected %v but got %v", expected, e.Error())
	}
}
func TestErrorfFormat(t *testing.T) {

	format := "Whatever %v"
	detail := "message"
	e := Errorf(format, detail)

	expected := utils.FmtPrefix + "/" + pkg + ": Whatever message"

	if e.Error() != expected {
		t.Errorf("Error is not well formatted, expected %s but got %s", expected, e.Error())
	}
}

func TestSprintFormat(t *testing.T) {

	msg := "Whatever message"
	e := Sprintf(msg)
	expected := utils.FmtPrefix + "/" + pkg + ": " + msg

	if e != expected {
		t.Errorf("String is not well formatted, expected %s but got %s", expected, e)
	}
}

func TestSprintfFormat(t *testing.T) {

	format := "Whatever %v"
	detail := "message"
	e := Sprintf(format, detail)

	expected := utils.FmtPrefix + "/" + pkg + ": Whatever message"

	if e != expected {
		t.Errorf("String is not well formatted, expected %s but got %s", expected, e)
	}
}
