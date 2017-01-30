package drivers_test

import (
	"github.com/spring1843/chat-server/drivers"
	"github.com/spring1843/chat-server/drivers/fake"
)

var _ drivers.Connection = fake.NewFakeConnection()
