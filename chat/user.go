package chat

import (
	"sync"

	"github.com/spring1843/chat-server/drivers"
)

// User is temporarily in connected to a chat server, and can be in certain channels
type User struct {
	conn drivers.Connection

	nickName     string
	lockNickName *sync.Mutex

	channel     string
	lockChannel *sync.Mutex

	ignoreList map[string]bool
	lockIgnore *sync.Mutex

	incoming chan string
	outgoing chan string
}

// NewUser returns a new new User
func NewUser(nickName string) *User {
	return &User{
		nickName:     nickName,
		channel:      "",
		ignoreList:   make(map[string]bool),
		incoming:     make(chan string),
		outgoing:     make(chan string),
		lockNickName: new(sync.Mutex),
		lockChannel:  new(sync.Mutex),
		lockIgnore:   new(sync.Mutex),
	}
}

// GetNickName returns the nickname of this user
func (u *User) GetNickName() string {
	u.lockNickName.Lock()
	defer u.lockNickName.Unlock()
	return u.nickName
}

// SetNickName sets the nickname for this user
func (u *User) SetNickName(nickName string) {
	u.lockNickName.Lock()
	defer u.lockNickName.Unlock()
	u.nickName = nickName
}

// GetChannel gets the current channel name for the user
func (u *User) GetChannel() string {
	u.lockChannel.Lock()
	defer u.lockChannel.Unlock()
	return u.channel
}

// SetChannel sets the current channel name for the user
func (u *User) SetChannel(name string) {
	u.lockChannel.Lock()
	defer u.lockChannel.Unlock()
	u.channel = name
}

// Ignore a user
func (u *User) Ignore(nickName string) {
	u.lockIgnore.Lock()
	defer u.lockIgnore.Unlock()
	u.ignoreList[nickName] = true
}

// HasIgnored checks to see if a user has ignored another user or not
func (u *User) HasIgnored(nickName string) bool {
	u.lockIgnore.Lock()
	defer u.lockIgnore.Unlock()
	if _, ok := u.ignoreList[nickName]; ok {
		return true
	}
	return false
}
