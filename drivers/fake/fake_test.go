package fake_test

import (
	"reflect"
	"testing"

	"github.com/spring1843/chat-server/drivers/fake"
)

func TestStringReadAndWrite(t *testing.T) {
	conn := fake.NewFakeConnection()
	msg1 := "foo"
	conn.WriteString(msg1)
	outgoing, err := conn.ReadString(len(msg1))
	if err != nil {
		t.Fatalf("Error reading from connection. Error: %s", err)
	}

	if outgoing != msg1 {
		t.Fatalf("Couldn't write and read to outgoing. Expected %q got %q.", msg1, outgoing)
	}
}

func TestReadAndWrite(t *testing.T) {
	conn := fake.NewFakeConnection()

	const msg = "foo"

	msg1 := []byte(msg)
	msg2 := make([]byte, len(msg), len(msg))

	conn.Write(msg1)

	n, err := conn.Read(msg2)
	if err != nil {
		t.Fatalf("Couldnt read from connection. Error %s.", err)
	}
	if n != len(msg) {
		t.Fatalf("Length of input and output do not match. Expected %d got %d", len(msg), n)
	}

	if !reflect.DeepEqual(msg1, msg2) {
		t.Fatalf("Input and output do not match, expected %q got %q.", string(msg1), string(msg2))
	}
}
