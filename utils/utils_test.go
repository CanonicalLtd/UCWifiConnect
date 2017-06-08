package utils

import (
	"os"
	"testing"
)

func TestHashFile(t *testing.T) {
	path := "./hash"
	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
	}
	SetHashFile(path)
	if HashFile != path {
		t.Errorf("Error. SetHashFile() failed to set HashFile var to %s. is still %s", path, HashFile)
	}
}

func TestHashIt(t *testing.T) {
	path := "/tmp/hash"
	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
	}
	SetHashFile(path)
	pwd := "12341234"
	_, e1 := HashIt(pwd)
	if e1 != nil {
		t.Errorf("Error %v HashIt returned an error", e1.Error())
	}
}

func TestMatchingHash(t *testing.T) {
	path := "/tmp/hash"
	_, err := os.Stat(path)
	if err != nil {
		os.Remove(path)
	}
	SetHashFile(path)
	pwd := "12341234"
	HashIt(pwd)
	match, e1 := MatchingHash(pwd)
	if e1 != nil {
		t.Errorf("Error %v MatchingHash returned an error", e1.Error())
	}
	if !match {
		t.Error("Error MatchingHash should have matched but did not")
	}
}
