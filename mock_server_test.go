package actioncable

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
)

const (
	fakeEndpoint = "ws://foo/ws"
)

// TestDialer merely makes sure that the wstest dialer is working correctly
func TestDialer(t *testing.T) {
	mockServer := &MockRailsServer{T: t}

	dialer := wstest.NewDialer(mockServer)
	dialer.HandshakeTimeout = time.Second * 2

	_, _, err := dialer.Dial(fakeEndpoint, nil)
	if err != nil {
		t.Fatal(err)
	}
}

// TestMockOnHeartbeat MUST start with ws://foo/ws, it tests that the client can receive
// three quick pings from a actioncable protocol following server
func TestMockOnHeartbeat(t *testing.T) {
	mockServer := &MockRailsServer{T: t, Behaviour: MockHeartbeat}

	dialer := wstest.NewDialer(mockServer)
	dialer.HandshakeTimeout = time.Second * 2

	client := NewClient(fakeEndpoint).WithDialer(dialer)

	called := make(chan struct{})
	count := 0

	client.OnHeartbeat(func(conn *websocket.Conn, payload *Payload) error {
		count++
		if count >= 4 {
			called <- struct{}{}
		}
		return nil
	})
	err := client.Serve()

	if err != nil {
		t.Fatal(err)
	}

	receiveSleepMs(2000, called, t)
}

// TestMockOnEvent tests the ability to setup an OnConnect callback
func TestMockOnConnect(t *testing.T) {
	mockServer := &MockRailsServer{T: t, Behaviour: MockConnect}

	dialer := wstest.NewDialer(mockServer)
	dialer.HandshakeTimeout = time.Second * 2

	client := NewClient(fakeEndpoint).WithDialer(dialer)

	called := make(chan struct{})

	client.OnConnect(func(conn *websocket.Conn) error {
		called <- struct{}{}
		return nil
	})
	err := client.Serve()

	if err != nil {
		t.Fatal(err)
	}

	receiveSleepMs(2000, called, t)
}

// TestMockOnEvent tests the ability to setup an OnEvent callback
func TestMockOnEvent(t *testing.T) {
	mockServer := &MockRailsServer{T: t, Behaviour: MockEvent}

	dialer := wstest.NewDialer(mockServer)
	dialer.HandshakeTimeout = time.Second * 2

	client := NewClient(fakeEndpoint).WithDialer(dialer)

	called := make(chan struct{})

	client.OnEvent("AgentChannel", func(conn *websocket.Conn, payload *Payload, error error) {
		called <- struct{}{}
		return
	})

	err := client.Serve()

	if err != nil {
		t.Fatal(err)
	}

	receiveSleepMs(2000, called, t)
}

func TestMockOnDisconnect(t *testing.T) {
	mockServer := &MockRailsServer{T: t, Behaviour: MockDisconnect}

	dialer := wstest.NewDialer(mockServer)
	dialer.HandshakeTimeout = time.Second * 1

	client := NewClient(fakeEndpoint).WithDialer(dialer)

	called := make(chan struct{})

	client.OnDisconnect(func(conn *websocket.Conn) error {
		called <- struct{}{}
		return nil
	})

	err := client.Serve()

	client.Close()

	if err != nil {
		t.Fatal(err)
	}

	receiveSleepMs(200, called, t)
}
