package errs_test

import (
	"testing"

	"github.com/spring1843/chat-server/src/shared/errs"
)

func TestCanFindErrorCause(t *testing.T) {
	cause := errs.Newf("Firs%s", "t")
	err := errs.Wrap(cause, "second")
	err = errs.Wrapf(err, "thir%s", "d")
	err = errs.Wrap(err, "last")

	if errs.Cause(err) != cause {
		t.Fatalf("Couldn't find cause, expected %s, got %s", cause, errs.Cause(err))
	}

	if errs.New("First").Error() != cause.Error() {
		t.Fatalf("Formatting doesn't match, expected %s, got %s", errs.New("first").Error(), cause.Error())
	}
}
