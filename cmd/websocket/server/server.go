package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() *Server {
	ch := make(chan *Server)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		ch <- &Server{conn: conn}
	})

	fmt.Println("Server is listening on port 8080")
	go log.Fatal(http.ListenAndServe(":8080", nil))
	return <-ch
}

func (s *Server) Test() {
	for {
		messageType, message, err := s.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("Received message: %s\n", message)

		if err := s.conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *Server) ReadMessage() ([]byte, error) {
	_, message, err := s.conn.ReadMessage()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Received message: %s\n", message)

	return message, err
}

func (s *Server) WriteMessage(message []byte) error {
	if err := s.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}
