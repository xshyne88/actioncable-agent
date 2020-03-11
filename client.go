package actioncable

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client for ActionCable Protocol websockets
type Client struct {
	mux       sync.Mutex
	address   string
	dialer    *websocket.Dialer
	conn      *websocket.Conn
	broadcast *Broadcast
	handlers  []func()
}

// Event ActionCable Protocol
type Event struct {
	Type string `json:"type"`

	Message    json.RawMessage `json:"message"`
	Data       json.RawMessage `json:"data"`
	Identifier *Identifier     `json:"identifier"`
}

// Payload gets returned from all of the callbacks
type Payload struct {
	event Event
}

const (
	welcome = "welcome"
	ping    = "ping"
)

// NewClient creates new Client Object
func NewClient(address string) *Client {
	return &Client{
		address:   address,
		broadcast: NewBroadcast(),
		handlers:  []func(){},
		dialer:    &websocket.Dialer{HandshakeTimeout: time.Second * 1},
	}
}

// WithDialer allows a custom dialer to be passed
func (c *Client) WithDialer(dialer *websocket.Dialer) *Client {
	c.dialer = dialer
	return c
}

// Serve dials back home and sets up the handlers
// note that conn MUST be set on the client before ranging over the handlers
func (c *Client) Serve() error {
	conn, _, err := c.dialer.Dial(c.address, nil)
	if err != nil {
		return err
	}

	c.conn = conn

	for _, h := range c.handlers {
		go h()
	}

	go readLoop(c)

	return nil
}

// loop pulls msgs off the stream and sends them to broadcast
func readLoop(c *Client) {
	evt := &Event{}
	for {
		if err := c.conn.ReadJSON(&evt); err != nil {
			c.broadcast.Error(err)
		}
		switch evt.Type {
		case welcome:
			c.broadcast.Connect(*evt)
		case ping:
			c.broadcast.Ping(*evt)
		default:
			c.broadcast.Event(*evt)
		}
	}
}

// OnConnect callback fired when Ws is first dialed
func (c *Client) OnConnect(cb func(conn *websocket.Conn) error) {
	f := func() {
		for {
			select {
			case <-c.broadcast.ConnectChan:
				cb(c.conn)
			case <-c.broadcast.DoneChan:
				return
			}
		}
	}
	c.handlers = append(c.handlers, f)
}

// OnEvent is fired when any message at all is recieved from the api
func (c *Client) OnEvent(chanName string, cb func(conn *websocket.Conn, payload *Payload, err error)) {
	f := func() {
		err := c.Subscribe(chanName)
		if err != nil {
			cb(c.conn, &Payload{}, err)
		}
		for {
			select {
			case event := <-c.broadcast.EventChan:
				if event.Identifier.Channel == chanName {
					cb(c.conn, &Payload{event: event}, nil)
				}
			case <-c.broadcast.DoneChan:
				return
			}
		}
	}
	c.handlers = append(c.handlers, f)
}

// OnHeartbeat provides a way to setup a handler on a ping message
func (c *Client) OnHeartbeat(cb func(conn *websocket.Conn, payload *Payload) error) {
	f := func() {
		for {
			select {
			case event := <-c.broadcast.HeartbeatChan:
				cb(c.conn, &Payload{event: event})
			case <-c.broadcast.DoneChan:
				return
			}
		}
	}
	c.handlers = append(c.handlers, f)
}

// OnDisconnect is the cb fired when the client disconnects
func (c *Client) OnDisconnect(cb func(conn *websocket.Conn) error) {
	f := func() {
		<-c.broadcast.DoneChan
		cb(c.conn)
	}
	c.handlers = append(c.handlers, f)
}

// Subscribe listens on a different channel
func (c *Client) Subscribe(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	cmd := NewSubscription(name)
	err := c.conn.WriteJSON(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Unsubscribe cancels subscriptions on a channel
func (c *Client) Unsubscribe(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	cmd := CancelSubscription(name)
	err := c.conn.WriteJSON(cmd)
	if err != nil {
		return err
	}
	return nil
}

// Close closes the client
func (c *Client) Close() {
	go func() {
		c.broadcast.Close()
	}()
}
