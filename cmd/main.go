package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var newID = 1

var currentConnections map[int]*websocket.Conn
var chatMessages [][]byte

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

func init() {
	currentConnections = make(map[int]*websocket.Conn)
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Обновляем соединение до веб-сокета
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error upgrading to websocket:", err)
		return
	}

	go handleConnection(conn)
}

func handleConnection(conn *websocket.Conn) {
	// Добавляем новое подключение в мапу текущих подключений
	curConnId := newID
	currentConnections[curConnId] = conn
	newID++

	msgChan := make(chan []byte)

	// При новом подключении покажем пользователю все старые сообщения
	for _, msg := range chatMessages {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println("error loading previous messages", err)
			return
		}
	}

	go func() {
		for {
			message, ok := <-msgChan
			if !ok {
				return
			}
			// Делаем что-нибудь с полученным сообщением из канала
			// Добавляем сообщение в текущие сообщения канала
			chatMessages = append(chatMessages, message)
			// Транслируем сообщения всем подключенным пользователям
			for id, user := range currentConnections {
				if err := user.WriteMessage(websocket.TextMessage, message); err != nil {
					log.Printf("error writing message to user %d %s", id, err)
				}
			}
		}
	}()

	// Горутина, которая запускает бесконечный цикл
	// для чтения из коннекта и записи результата в канал
	go func() {
		for {
			defer conn.Close()
			defer delete(currentConnections, curConnId)

			_, message, err := conn.ReadMessage()
			if err != nil {
				close(msgChan)
				break
			}
			msgChan <- message
		}
	}()
}

func main() {
	http.HandleFunc("/ws", handler)
	http.ListenAndServe(":8080", nil)
}
