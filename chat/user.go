package chat

import "sync"

// User is temporarily in connected to a chat server, and can be in certain channels
type User struct {
	conn       Connection
	nickName   string
	channel    string
	ignoreList map[string]bool
	incoming   chan string
	outgoing   chan string
	lock       *sync.Mutex
}

// NewUser returns a new new User
func NewUser(nickName string) *User {
	return &User{
		nickName:   nickName,
		channel:    "",
		ignoreList: make(map[string]bool),
		incoming:   make(chan string),
		outgoing:   make(chan string),
		lock:       new(sync.Mutex),
	}
}

// GetChannel gets the current channel name for the user
func (u *User) GetChannel() string {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.channel
}

// GetNickName returns the nickname of this user
func (u *User) GetNickName() string {
	u.lock.Lock()
	defer u.lock.Unlock()
	return u.nickName
}

// SetNickName sets the nickname for this user
func (u *User) SetNickName(nickName string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.nickName = nickName
}

// SetChannel sets the current channel name for the user
func (u *User) SetChannel(name string) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.channel = name
}

// Ignore a user
func (u *User) Ignore(nickName string) {
	u.ignoreList[nickName] = true
}

// HasIgnored checks to see if a user has ignored another user or not
func (u *User) HasIgnored(nickName string) bool {
	if _, ok := u.ignoreList[nickName]; ok {
		return true
	}
	return false
}
