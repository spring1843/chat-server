package drivers_test

import (
	"github.com/spring1843/chat-server/src/drivers"
	"github.com/spring1843/chat-server/src/drivers/fake"
)

var _ drivers.Connection = fake.NewFakeConnection()
