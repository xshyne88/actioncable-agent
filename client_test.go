package client

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func assert(act, exp interface{}, t *testing.T) {
	if act != exp {
		t.Errorf("Act: %v: Exp: %v", act, exp)
	}
}

func assertDeep(act, exp interface{}, t *testing.T) {
	if !reflect.DeepEqual(exp, act) {
		t.Errorf("Act: %v: Exp: %v", act, exp)
	}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Error(err)
	}
}

func sleepMs(t time.Duration) {
	time.Sleep(time.Millisecond * t)
}

func TestClientOnConnect(t *testing.T) {
	dialer := &websocket.Dialer{HandshakeTimeout: time.Second * 1}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	checkError(err, t)

	called := make(chan struct{})
	client.OnConnect(func(conn *websocket.Conn) error {
		called <- struct{}{}
		return nil
	})

	err = client.Serve()
	checkError(err, t)

	sleepMs(500)

	assert(<-called, true, t)
}

func TestClientOnDisconnect(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	checkError(err, t)

	var called bool
	var count int

	client.OnDisconnect(func(conn *websocket.Conn) error {
		called = true
		count++
		return nil
	})

	err = client.Serve()
	checkError(err, t)

	sleepMs(200)
	client.Close()

	// give time for cb
	sleepMs(200)
	assert(called, true, t)
	assert(count, 1, t)
}

// go test *.go -v -run TestClientOnHeartbeat -count=1
func TestClientOnHeartbeat(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	if err != nil {
		t.Error(err)
	}

	var actual Payload
	done := make(chan struct{})

	client.OnHeartbeat(func(conn *websocket.Conn, sample *Payload) error {
		actual.event = sample.event
		done <- struct{}{}
		return nil
	})

	err = client.Serve()
	checkError(err, t)

	// give enough time for ping (3500 wasnt enough)
	receiveSleep(4000, done)
	assert(actual.event.Type, "ping", t)
}

func receiveSleep(t time.Duration, done chan struct{}) {
	select {
	case <-time.After(t * time.Millisecond):
		return
	case <-done:
		return
	}
}

func TestClientOnEvent(t *testing.T) {
	dialer := &websocket.Dialer{
		HandshakeTimeout: time.Second * 1,
	}
	client, err := NewClient("wss://api.rmm.dev/cable", dialer)
	if err != nil {
		t.Error(err)
	}

	err = client.OnEvent("agent", func(conn *websocket.Conn, sample *Payload) error {
		var str string
		err := conn.ReadJSON(&str)
		if err != nil {
			panic(err)
		}
		fmt.Print(str)
		return nil
	})

	client.Serve()
	err = client.Serve()
	checkError(err, t)

	sleepMs(2000)
}

func TestEnd2End(t *testing.T) {

}
