package client

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xshyne88/agent/event"
)

// PingMessage mocks heartbeats from rails in a faster time 100ms
type PingMessage struct {
	time int64
	Type string
}

// MockRailsServer is a struct for helping build the mock server
type MockRailsServer struct {
	Upgrader websocket.Upgrader
	Done     chan struct{}
	T        *testing.T
}

// Implements http.Handler interface, sends welcome and 4 pings before exiting
func (s *MockRailsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.T.Helper()
	var conn *websocket.Conn

	checkWsPath(s.T, r.URL.Path)

	s.Done = make(chan struct{})
	defer close(s.Done)

	conn, _ = s.Upgrader.Upgrade(w, r, nil) // should NEVER error since its a wstest dependency

	wel := event.Event{
		Type: "welcome",
	}
	conn.WriteJSON(wel)

	sentPings := 0
	for {
		t := time.NewTicker(time.Duration(time.Millisecond * 100))
		if sentPings > 3 {
			return
		}
		select {
		case <-t.C:
			sendPing(conn)
			sentPings++
		case <-s.Done:
			t.Stop()
			return
		}
	}
}

// sendPing is a helper to send a ping
func sendPing(conn *websocket.Conn) {
	p := &PingMessage{
		time: time.Now().Unix(),
		Type: "ping",
	}
	conn.WriteJSON(p)
}

// checkWsPath checks the absolute url of the websocket path and errors if it doesnt start with ws
func checkWsPath(t *testing.T, path string) {
	if path != "/ws" {
		t.Errorf("Path does not contain ws:// or wss://")
	}
}