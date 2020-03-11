package actioncable

// Broadcast is a hub for broadcasting messages
type Broadcast struct {
	EventChan     chan Event
	ConnectChan   chan Event
	HeartbeatChan chan Event

	ErrorChan chan error

	DoneChan chan struct{}
}

// NewBroadcast is a constructor for Broadcast
func NewBroadcast() *Broadcast {
	return &Broadcast{
		EventChan:     make(chan Event),
		ConnectChan:   make(chan Event),
		HeartbeatChan: make(chan Event),
		DoneChan:      make(chan struct{}),
		ErrorChan:     make(chan error),
	}
}

// Error puts error on Error channel
func (b *Broadcast) Error(err error) {
	go func() {
		b.ErrorChan <- err
		return
	}()
}

// Event puts event on event channel
func (b *Broadcast) Event(e Event) {
	go func() {
		b.EventChan <- e
		return
	}()
}

// Ping channel
func (b *Broadcast) Ping(e Event) {
	go func() {
		b.HeartbeatChan <- e
		return
	}()
}

// Connect puts event on connect channel
func (b *Broadcast) Connect(e Event) {
	go func() {
		b.ConnectChan <- e
		return
	}()
}

// Close puts struct onto done channel
func (b *Broadcast) Close() {
	b.DoneChan <- struct{}{}
}
