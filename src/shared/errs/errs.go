package errs

import "github.com/pkg/errors"

// New returns an error with the supplied message.
func New(message string) error {
	return errors.New(message)
}

// Newf formats according to a format specifier and returns the string
func Newf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// Wrap If err is nil, Wrap returns nil.
func Wrap(err error, message string) error {
	return errors.Wrap(err, message)
}

// Wrapf returns an error annotating err with the format specifier.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Cause returns the underlying cause of the error, if possible.
func Cause(err error) error {
	return errors.Cause(err)
}
