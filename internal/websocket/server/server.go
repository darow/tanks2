package server

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"
	"sync"

	"myebiten/internal/models"

	"github.com/gorilla/websocket"
)

type InputStore struct {
	sync.Mutex
	input models.Input
}

type Server struct {
	charInputConn    *websocket.Conn
	thingsUpdateConn *websocket.Conn
	mapUpdateConn    *websocket.Conn
	inputStore       *InputStore
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New(port string) *Server {
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
		addr := ":" + port
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	s := &Server{
		charInputConn:    <-ch1,
		thingsUpdateConn: <-ch2,
		mapUpdateConn:    <-ch3,
		inputStore:       &InputStore{},
	}

	go s.ReceiveUpdates()
	log.Println("client connected")

	return s
}

func (s *Server) ReceiveUpdates() {
	for {
		_, rawMessage, err := s.charInputConn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
			log.Fatal()
		}

		var input models.Input
		if err := json.Unmarshal(rawMessage, &input); err != nil {
			continue
		}

		// Preserve Shoot: if new input has Shoot=false but previous had Shoot=true,
		// keep Shoot=true so the shoot gets processed
		oldShoot := s.inputStore.input.Shoot
		newShoot := input.Shoot
		input.Shoot = oldShoot || newShoot

		s.inputStore.Lock()
		s.inputStore.input = input
		s.inputStore.Unlock()
	}
}

func (s *Server) GetInput() models.Input {
	s.inputStore.Lock()
	defer s.inputStore.Unlock()
	return s.inputStore.input
}

func (s *Server) SetInputShootFalse() {
	if s == nil || s.inputStore == nil {
		return
	}

	s.inputStore.Lock()
	defer s.inputStore.Unlock()
	s.inputStore.input.Shoot = false
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
