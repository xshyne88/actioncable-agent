package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xshyne88/agent/command"
)

type Client struct {
	address   string
	conn      *websocket.Conn
	broadcast *Broadcast
}

type Event struct {
	Type string `json:"type"`

	Message    json.RawMessage     `json:"message"`
	Data       json.RawMessage     `json:"data"`
	Identifier *command.Identifier `json:"identifier"`
}

// NewClient creates new Client Object
func NewClient(address string, dialer *websocket.Dialer) (*Client, error) {
	conn, _, err := dialer.Dial(address, make(http.Header))

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	c := &Client{
		address:   address,
		conn:      conn,
		broadcast: NewBroadcast(),
	}

	go runLoop(c)

	return c, nil
}

func runLoop(c *Client) {
	getWsMessages(c.conn, c)
}

func getWsMessages(conn *websocket.Conn, c *Client) {
	evt := &Event{}
	for {
		if err := conn.ReadJSON(&evt); err != nil {
			c.broadcast.Error(err)
		}
		if evt.Type == "welcome" {
			c.broadcast.Connect(*evt)
		} else {
			c.broadcast.Event(*evt)
		}
	}
}

// OnConnect is
func (c *Client) OnConnect(cb func(conn *websocket.Conn) error) error {
	go func() {
		for {
			select {
			case <-c.broadcast.ConnectChan:
				cb(c.conn)
			case <-c.broadcast.DoneChan:
				return
			default:
			}
		}
	}()
	return nil
}

// Payload gets returned to the callback user
type Payload struct {
	event Event
}

// OnEvent is
func (c *Client) OnEvent(eventType string, cb func(conn *websocket.Conn, payload *Payload) error) error {
	go func() {
		for {
			select {
			case event := <-c.broadcast.EventChan:
				switch event.Type {
				case eventType:
					cb(c.conn, &Payload{event: event})
				default:
					fmt.Printf("non %s message received", eventType)
				}
			case <-c.broadcast.DoneChan:
				return
			default:
			}
		}
	}()
	return nil
}

// OnDisconnect is
func (c *Client) OnDisconnect(cb func(conn *websocket.Conn) error) error {
	// Todo: needs to make sure it runs before close -
	go func() {
		for {
			select {
			case <-c.broadcast.DoneChan:
				cb(c.conn)
				return
			}
		}
	}()
	return nil
}

// Close closes the client
func (c *Client) Close() {
	go func() {
		c.broadcast.Close()
	}()
}
