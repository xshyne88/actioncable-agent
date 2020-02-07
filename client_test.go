package client

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func assert(act, exp interface{}, t *testing.T) {
	if act != exp {
		t.Errorf("Act: %v: Exp: %v", act, exp)
	}
}

func TestClientOnConnect(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	if err != nil {
		t.Error(err)
	}

	called := false
	client.OnConnect(func(conn *websocket.Conn) error {
		called = true
		return nil
	})

	time.Sleep(time.Second * 2)

	assert(called, true, t)
}

func TestClientOnDisconnect(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	if err != nil {
		t.Error(err)
	}

	count := 0
	called := false
	client.OnDisconnect(func(conn *websocket.Conn) error {
		count++
		called = true
		return nil
	})

	client.Close()

	time.Sleep(time.Second * 2)

	t.Logf("count: %d", count)
	assert(called, true, t)
}

func TestClientOnEvent(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	if err != nil {
		t.Error(err)
	}

	var actual Payload
	client.OnEvent("ping", func(conn *websocket.Conn, sample *Payload) error {
		actual.event = sample.event
		t.Logf("event: %+v", sample)
		return nil
	})

	time.Sleep(time.Second * 2)

	assert(actual.event.Type, "ping", t)
}
