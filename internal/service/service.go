package service

import (
	"strconv"
	"sync"

	"github.com/dkshi/chopchop"
	"github.com/gorilla/websocket"
)

type Service struct {
	users     map[int]*chopchop.User
	companies map[int]int
	currID    int

	mu *sync.Mutex
}

func NewService() *Service {
	return &Service{
		currID:    1,
		users:     make(map[int]*chopchop.User),
		companies: make(map[int]int),
		mu:        &sync.Mutex{},
	}
}

func (s *Service) AddConn(conn *websocket.Conn) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	newID := s.currID
	s.users[newID] = chopchop.NewUser(newID, "<noname>", conn)
	s.currID++

	return newID
}

func (s *Service) SendMessageCompany(msg []byte, connID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if receiverID, ok := s.companies[connID]; ok {
		newMessage := []byte(s.users[connID].Username + ": " + string(msg))
		s.users[connID].Conn.WriteMessage(websocket.TextMessage, newMessage)
		s.users[receiverID].Conn.WriteMessage(websocket.TextMessage, newMessage)
	} else {
		s.users[connID].Conn.WriteMessage(websocket.TextMessage, []byte("you don't have a company!"))
	}
	
}

func (s *Service) BroadcastMessage(msg []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, user := range s.users {
		if err := user.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) MakeCompany(connID int, stringReceiverID string) *chopchop.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	receiverID, err := strconv.Atoi(stringReceiverID)
	if err != nil {
		return chopchop.NewError(666, err.Error())
	}

	if connID == receiverID {
		return chopchop.NewError(666, "error: you cannot make company with yourself")
	}

	if _, ok := s.users[receiverID]; !ok {
		return chopchop.NewError(666, "error: there are no such user with id: "+strconv.Itoa(receiverID))
	}

	if _, ok := s.companies[connID]; ok {
		return chopchop.NewError(666, "error: you are already in company")
	}

	if _, ok := s.companies[receiverID]; ok {
		return chopchop.NewError(666, "error: user id: "+strconv.Itoa(receiverID)+" are already in company")
	}

	s.companies[connID] = receiverID
	s.companies[receiverID] = connID
	return chopchop.NewError(0, "successfully made company with user id: "+strconv.Itoa(receiverID))
}

func (s *Service) BreakCompany(connID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.companies, s.companies[connID])
	delete(s.companies, connID)
}

func (s *Service) DeleteConn(connID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.BreakCompany(connID)
	delete(s.users, connID)
}

func (s *Service) WriteConns(connID int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	reply := ""
	for id, user := range s.users {
		line := strconv.Itoa(id) + " " + user.Username
		if id == connID {
			line += " (you)"
		}
		reply += line + "\n"
	}
	s.users[connID].Conn.WriteMessage(websocket.TextMessage, []byte(reply))
}

func (s *Service) RenameConn(connID int, newName string) *chopchop.Error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(newName) > 24 {
		return chopchop.NewError(666, "error: your new name is too long")
	}

	s.users[connID].Rename(newName)
	return chopchop.NewError(0, "name was changed succesfully")
}
