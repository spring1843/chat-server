package main

import (
	"testing"
)

func Test_CommandLineArguments(t *testing.T) {
	validCmd := []string{`chatserver`, `-c`, `somefile`}
	invalidCmd1 := []string{`chatserver`}
	invalidCmd2 := []string{`chatserver`, `-c`}
	invalidCmd3 := []string{`chatserver`, `-d`}

	if validateCommandArguments(validCmd) != true {
		t.Errorf("Valid arguments were determined invalid")
	}

	if validateCommandArguments(invalidCmd1) != false {
		t.Errorf("Invalid arguments were determined valid")
	}

	if validateCommandArguments(invalidCmd2) != false {
		t.Errorf("Invalid arguments were determined valid")
	}

	if validateCommandArguments(invalidCmd3) != false {
		t.Errorf("Invalid arguments were determined valid")
	}
}
