package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"sync"

	"myebiten/internal/models"

	"github.com/gorilla/websocket"
)

type InputStore struct {
	sync.Mutex
	inputs map[int]models.Input
}

type Server struct {
	charInputConns    map[int]*websocket.Conn
	thingsUpdateConns map[int]*websocket.Conn
	mapUpdateConns    map[int]*websocket.Conn
	inputStore        *InputStore
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type playerConn struct {
	playerID int
	conn     *websocket.Conn
}

func New(port string, playersCount int) *Server {
	ch1 := make(chan playerConn)
	ch2 := make(chan playerConn)
	ch3 := make(chan playerConn)

	http.HandleFunc("/ws1", connectionHandler(playersCount, ch1))
	http.HandleFunc("/ws2", connectionHandler(playersCount, ch2))
	http.HandleFunc("/ws3", connectionHandler(playersCount, ch3))
	http.HandleFunc("/players_count", playersCountHandler(playersCount))

	go func() {
		log.Printf("Server is listening on port %s\n", port)
		addr := ":" + port
		log.Fatal(http.ListenAndServe(addr, nil))
	}()

	inputConns, thingsUpdateConns, mapUpdateConns := waitForClientConnections(playersCount, ch1, ch2, ch3)

	s := &Server{
		charInputConns:    inputConns,
		thingsUpdateConns: thingsUpdateConns,
		mapUpdateConns:    mapUpdateConns,
		inputStore:        &InputStore{inputs: map[int]models.Input{}},
	}

	for playerID, conn := range s.charInputConns {
		go s.ReceiveUpdates(playerID, conn)
	}
	log.Printf("%d clients connected\n", playersCount-1)

	return s
}

func playersCountHandler(playersCount int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(playersCount); err != nil {
			log.Println(err)
		}
	}
}

func connectionHandler(playersCount int, ch chan<- playerConn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		playerID, err := strconv.Atoi(r.URL.Query().Get("player_id"))
		if err != nil || playerID <= 0 || playerID >= playersCount {
			http.Error(w, fmt.Sprintf("player_id must be from 1 to %d", playersCount-1), http.StatusBadRequest)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		ch <- playerConn{playerID: playerID, conn: conn}
	}
}

func waitForClientConnections(
	playersCount int,
	inputCh <-chan playerConn,
	thingsUpdateCh <-chan playerConn,
	mapUpdateCh <-chan playerConn,
) (map[int]*websocket.Conn, map[int]*websocket.Conn, map[int]*websocket.Conn) {
	remotePlayersCount := playersCount - 1
	inputConns := make(map[int]*websocket.Conn, remotePlayersCount)
	thingsUpdateConns := make(map[int]*websocket.Conn, remotePlayersCount)
	mapUpdateConns := make(map[int]*websocket.Conn, remotePlayersCount)

	for len(inputConns) < remotePlayersCount ||
		len(thingsUpdateConns) < remotePlayersCount ||
		len(mapUpdateConns) < remotePlayersCount {
		select {
		case pc := <-inputCh:
			inputConns[pc.playerID] = pc.conn
		case pc := <-thingsUpdateCh:
			thingsUpdateConns[pc.playerID] = pc.conn
		case pc := <-mapUpdateCh:
			mapUpdateConns[pc.playerID] = pc.conn
		}
	}

	return inputConns, thingsUpdateConns, mapUpdateConns
}

func (s *Server) ReceiveUpdates(playerID int, conn *websocket.Conn) {
	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			log.Println(runtime.Caller(1))
			log.Println(err)
			log.Fatal()
		}

		var input models.Input
		if err := json.Unmarshal(rawMessage, &input); err != nil {
			continue
		}

		s.inputStore.Lock()
		s.inputStore.inputs[playerID] = input
		s.inputStore.Unlock()
	}
}

func (s *Server) GetInput(playerID int) models.Input {
	s.inputStore.Lock()
	defer s.inputStore.Unlock()
	return s.inputStore.inputs[playerID]
}

func (s *Server) WriteThingsMessage(message []byte) error {
	for _, conn := range s.thingsUpdateConns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) WriteMapMessage(message []byte) error {
	for _, conn := range s.mapUpdateConns {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return err
		}
	}

	return nil
}
