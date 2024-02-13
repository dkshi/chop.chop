package repository

import (
	"github.com/dkshi/chopchop"
	"github.com/gorilla/websocket"
)

type ChatInterface interface {
	GetUser(key int) (*chopchop.User, bool)
	AddUser(conn *websocket.Conn) int
	GetCompany(key int) (int, bool)
	GetUsers() []*chopchop.User
	AddCompany(key, value int)
	DeleteCompany(key int)
	DeleteUser(key int)
}

type Repository struct {
	ChatInterface
}

func NewRepository() *Repository {
	return &Repository{
		ChatInterface: NewChat(),
	}
}
