package fake

import "testing"

func TestCanIgnore(t *testing.T) {
	fakeReader := NewFakeConnection()

	fakeReader.incoming = []byte("foo\n")
}
