package repository

import (
	"sync"

	"github.com/dkshi/chopchop"
	"github.com/gorilla/websocket"
)

type Chat struct {
	users     sync.Map
	companies sync.Map
	newID     int
}

func NewChat() *Chat {
	return &Chat{
		newID: 1,
	}
}

func (c *Chat) GetUsers() []*chopchop.User {
	currUsers := make([]*chopchop.User, 0, c.newID)
	c.users.Range(func(key, value any) bool {
		currUsers = append(currUsers, value.(*chopchop.User))
		return true
	})
	return currUsers
}

func (c *Chat) GetUser(key int) (*chopchop.User, bool) {
	user, ok := c.users.Load(key)
	if !ok {
		return &chopchop.User{}, ok
	}
	return user.(*chopchop.User), ok
}

func (c *Chat) AddUser(conn *websocket.Conn) int {
	c.users.Store(c.newID, chopchop.NewUser(c.newID, "<noname>", conn))
	c.newID++
	return c.newID - 1
}

func (c *Chat) DeleteUser(key int) {
	c.users.Delete(key)
}

func (c *Chat) AddCompany(key, value int) {
	c.companies.Store(key, value)
}

func (c *Chat) GetCompany(key int) (int, bool) {
	company, ok := c.companies.Load(key)
	if !ok {
		return 0, ok
	}
	return company.(int), ok
}

func (c *Chat) DeleteCompany(key int) {
	c.companies.Delete(key)
}
