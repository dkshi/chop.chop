package chopchop

import "github.com/gorilla/websocket"

type User struct {
	ID       int
	Username string
	Conn     *websocket.Conn
}

func NewUser(id int, name string, conn *websocket.Conn) *User {
	return &User{
		ID:       id,
		Username: name,
		Conn:     conn,
	}
}