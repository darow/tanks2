package server

import (
	"log"
	"net/http"
	"runtime"
	"sync"

	"github.com/gorilla/websocket"
)

type MessageStore struct {
	sync.Mutex
	message []byte
}
type Server struct {
	charInputConn    *websocket.Conn
	thingsUpdateConn *websocket.Conn
	mapUpdateConn    *websocket.Conn
	msgStore         *MessageStore
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() *Server {
	ch1 := make(chan *websocket.Conn)
	ch2 := make(chan *websocket.Conn)
	ch3 := make(chan *websocket.Conn)

	http.HandleFunc("/ws1", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		ch1 <- conn
	})

	http.HandleFunc("/ws2", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		ch2 <- conn
	})

	http.HandleFunc("/ws3", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		ch3 <- conn
	})

	go func() {
		log.Println("Server is listening on port 8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	s := &Server{
		charInputConn:    <-ch1,
		thingsUpdateConn: <-ch2,
		mapUpdateConn:    <-ch3,
		msgStore:         &MessageStore{},
	}

	go s.ReceiveUpdates()
	log.Println("client connected")

	return s
}

func (s *Server) ReceiveUpdates() {
	for {
		_, message, err := s.charInputConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
			log.Fatal()
		}

		s.msgStore.Lock()
		s.msgStore.message = message
		s.msgStore.Unlock()
	}
}

func (s *Server) ReadMessage() []byte {
	s.msgStore.Lock()
	message := s.msgStore.message
	s.msgStore.Unlock()

	return message
}

func (s *Server) WriteThingsMessage(message []byte) error {
	if err := s.thingsUpdateConn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}

func (s *Server) WriteMapMessage(message []byte) error {
	if err := s.mapUpdateConn.WriteMessage(websocket.TextMessage, message); err != nil {
		return err
	}

	return nil
}
