package logs_test

import (
	"errors"
	"os"
	"os/exec"
	"testing"

	"github.com/spring1843/chat-server/src/shared/logs"
)

// TestFatalExitOnFatalError tests that Fatal failure ends the process with 1 status
func TestFatalExitOnFatalError(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		logs.FatalIfErrf(errors.New("Planned Fatal Error"), "Should fail with exit status %d", 1)
		return
	} else {
		logs.Infof("Should fail with exit status %d", 0)
		logs.WarnIfErrf(errors.New("Planned Warning"), "Should fail with exit status %d", 0)
		logs.ErrIfErrf(errors.New("Planned Error"), "Should fail with exit status %d", 0)
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalExitOnFatalError")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()

	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}

	t.Fatalf("process ran with err %s, want exit status 1", err)
}
