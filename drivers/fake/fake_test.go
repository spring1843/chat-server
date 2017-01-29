package fake_test

import (
	"testing"

	"github.com/spring1843/chat-server/drivers/fake"
	"reflect"
)

func TestIncomingAndOutGoing(t *testing.T) {
	conn := fake.NewFakeConnection()

	msg1,msg2:="foo","bar"

	conn.SetIncoming(msg1)
	incoming := conn.GetIncoming()
	if incoming!= msg1 {
		t.Fatalf("Couldn't write and read to incoming. Expected %q got %q.", msg1, incoming)
	}
	
	conn.SetOutgoing(msg2)
	outgoing := conn.GetOutgoing()
	if outgoing!= msg2 {
		t.Fatalf("Couldn't write and read to outgoing. Expected %q got %q.", msg2, outgoing)
	}
}

func TestReadAndWrite(t *testing.T) {
	conn := fake.NewFakeConnection()

	const msg = "foo"

	msg1 :=[]byte(msg)
	msg2 := make([]byte, len(msg), len(msg))

	conn.Write(msg1)

	n, err := conn.Read(msg2)
		if err!= nil{
		t.Fatalf("Couldnt read from connection. Error %s.", err)
	}
	if n != len(msg) {
		t.Fatalf("Length of input and output do not match. Expected %d got %d", len(msg), n)
	}

	if !reflect.DeepEqual(msg1, msg2) {
		t.Fatalf("Input and output do not match, expected %q got %q.", string(msg1), string(msg2))
	}
}
