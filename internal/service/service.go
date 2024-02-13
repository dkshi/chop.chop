package service

import (
	"strconv"

	"github.com/dkshi/chopchop/internal/repository"
	"github.com/gorilla/websocket"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) AddConn(conn *websocket.Conn) int {
	return s.repo.AddUser(conn)
}

// Отправить сообщение текущую компанию
func (s *Service) SendMessageCompany(msg []byte, connID int) {
	currUser, _ := s.repo.GetUser(connID)

	if receiverID, ok := s.repo.GetCompany(connID); ok {
		currReceiver, _ := s.repo.GetUser(receiverID)

		newMessage := []byte(currUser.Username + ": " + string(msg))

		currUser.Conn.WriteMessage(websocket.TextMessage, newMessage)
		currReceiver.Conn.WriteMessage(websocket.TextMessage, newMessage)
		return
	}

	currUser.Conn.WriteMessage(websocket.TextMessage, []byte("you don't have a company!"))
}

// Отправить сообщение всем пользователям
func (s *Service) BroadcastMessage(msg []byte) {
	users := s.repo.GetUsers()
	for _, user := range users {
		user.Conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// Создать компанию с пользователем
func (s *Service) MakeCompany(connID int, stringReceiverID string) {
	currUser, _ := s.repo.GetUser(connID)
	receiverID, err := strconv.Atoi(stringReceiverID)
	currReceiver, _ := s.repo.GetUser(receiverID)

	if err != nil {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("incorrect format of user id"))
		return
	}

	if connID == receiverID {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("error: you cannot make company with yourself"))
		return
	}

	if _, ok := s.repo.GetUser(receiverID); !ok {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("error: there are no such user with id: "+strconv.Itoa(receiverID)))
		return
	}

	if _, ok := s.repo.GetCompany(connID); ok {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("error: you are already in company"))
		return
	}

	if _, ok := s.repo.GetCompany(receiverID); ok {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("error: user id: "+strconv.Itoa(receiverID)+" are already in company"))
		return
	}

	s.repo.AddCompany(connID, receiverID)
	s.repo.AddCompany(receiverID, connID)

	successMsg := []byte("successfully made company with user id: " + strconv.Itoa(receiverID))
	currUser.Conn.WriteMessage(websocket.TextMessage, successMsg)
	currReceiver.Conn.WriteMessage(websocket.TextMessage, successMsg)
}

// Разорвать компанию с пользователем
func (s *Service) BreakCompany(connID int) {
	currCompany, _ := s.repo.GetCompany(connID)
	s.repo.DeleteCompany(currCompany)
	s.repo.DeleteCompany(connID)
}

func (s *Service) DeleteConn(connID int) {
	s.BreakCompany(connID)
	s.repo.DeleteUser(connID)
}

func (s *Service) WriteConns(connID int) {
	reply := ""
	currUser, _ := s.repo.GetUser(connID)
	currUsers := s.repo.GetUsers()
	for _, user := range currUsers {
		line := strconv.Itoa(user.ID) + " " + user.Username
		if user.ID == connID {
			line += " (you)"
		}
		reply += line + "\n"
	}
	currUser.Conn.WriteMessage(websocket.TextMessage, []byte(reply))
}

func (s *Service) RenameConn(connID int, newName string) {
	currUser, _ := s.repo.GetUser(connID)
	if len(newName) > 24 {
		currUser.Conn.WriteMessage(websocket.TextMessage, []byte("error: your new name is too long"))
		return
	}

	currUser.Username = newName
	currUser.Conn.WriteMessage(websocket.TextMessage, []byte("name was changed successfully"))
}
