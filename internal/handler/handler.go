package handler

import (
	"log"
	"net/http"

	"github.com/dkshi/chopchop/internal/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

type Handler struct {
	srv *service.Service
}

func NewHandler(srv *service.Service) *Handler {
	return &Handler{
		srv: srv,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/chat", h.ChatHandler)

	return r
}

func (h *Handler) ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Обновляем соединение до веб-сокета
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading to websocket:", err)
		return
	}

	go h.handleConnection(conn)
}

func (h *Handler) handleConnection(conn *websocket.Conn) {
	// Добавляем новое подключение в мапу текущих подключений
	connID := h.srv.AddConn(conn)

	msgChan := make(chan []byte)

	h.srv.WriteConns(connID)

	go func() {
		for {
			message, ok := <-msgChan
			if !ok {
				return
			}
			// Делаем что-нибудь с полученным сообщением из канала
			strMessage := string(message)
			if len(strMessage) >= 7 && strMessage[:7] == "rename " && len(strMessage[7:]) != 0 {
				h.srv.RenameConn(connID, strMessage[7:])
				continue
			}
			if len(strMessage) >= 8 && strMessage[:8] == "company " && len(strMessage[8:]) != 0 {
				h.srv.MakeCompany(connID, strMessage[8:])
				continue
			}
			if len(strMessage) >= 4 && strMessage[:4] == "list" && len(strMessage[4:]) == 0 {
				h.srv.WriteConns(connID)
				continue
			}
			if len(strMessage) >= 5 && strMessage[:5] == "leave" && len(strMessage[5:]) == 0 {
				h.srv.BreakCompany(connID)
				continue
			}
			if len(strMessage) >= 4 && strMessage[:4] == "help" && len(strMessage[4:]) == 0 {
				h.srv.WriteCmdList(connID)
				continue
			}
			h.srv.SendMessageCompany(message, connID)
		}
	}()

	// Горутина, которая запускает бесконечный цикл
	// для чтения из коннекта и записи результата в канал
	go func() {
		for {
			defer conn.Close()
			defer h.srv.DeleteConn(connID)

			_, message, err := conn.ReadMessage()
			if err != nil {
				close(msgChan)
				break
			}
			msgChan <- message
		}
	}()
}
