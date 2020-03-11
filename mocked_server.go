package actioncable

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// PingMessage mocks heartbeats from rails in a faster time 100ms
type PingMessage struct {
	time int64
	Type string
}

// Behaviour allows the consumer to pass in the type of behaviour
// he would like the server to execute,
type Behaviour int

const (
	// MockConnect mocks OnConnect's callback behaviour
	MockConnect Behaviour = iota

	// MockHeartbeat sends three pings
	MockHeartbeat

	// MockDisconnect mocks OnDisconnect's callback behaviour
	MockDisconnect

	// MockEvent mocks OnEvent's callback behaviour
	MockEvent
)

// MockRailsServer is a struct for helping build the mock server
type MockRailsServer struct {
	Upgrader  websocket.Upgrader
	Done      chan struct{}
	T         *testing.T
	Behaviour Behaviour
}

// Implements http.Handler interface, sends welcome and 4 pings before exiting
// note the upgrader should not ever error
func (s *MockRailsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.T.Helper()
	var conn *websocket.Conn

	checkWsPath(s.T, r.URL.Path)

	s.Done = make(chan struct{})
	defer close(s.Done)

	conn, _ = s.Upgrader.Upgrade(w, r, nil)

	wel := Event{
		Type: welcome,
	}
	conn.WriteJSON(wel)

	switch s.Behaviour {
	case MockHeartbeat:
		sendHeartbeats(conn, s)
	case MockDisconnect:
		s.Done <- struct{}{}
		sleepMs(10)
	case MockEvent:
		sendEventBehaviour(conn, s)
	case MockConnect:
		return
	}
}

func sendEventBehaviour(conn *websocket.Conn, s *MockRailsServer) {
	var cmd Command
	for {
		err := conn.ReadJSON(&cmd)
		if err != nil {
			panic(err)
		}
		if cmd.Identifier.Channel == "AgentChannel" {
			err := conn.WriteJSON(&Event{
				Type: welcome,
				Identifier: &Identifier{
					Channel: "AgentChannel",
				},
			})
			if err != nil {
				panic(err)
			}

			fakeBroadcastMessage := []byte(`{"foo": "bar"}`)

			err = conn.WriteJSON(&Event{
				Type: "event",
				Identifier: &Identifier{
					Channel: "AgentChannel",
				},
				Data: fakeBroadcastMessage,
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

// sendHeartbeats three total heartbeats every 100ms
func sendHeartbeats(conn *websocket.Conn, s *MockRailsServer) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	sentPings := 0

	for {
		<-ticker.C
		sendPing(conn)
		sentPings++
		if sentPings >= 4 {
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
		t.Fatal("Path does not contain ws:// or wss://")
	}
}
