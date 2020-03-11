package actioncable

import (
	"encoding/json"
)

// Command struct that follows ActionCable protocol
type Command struct {
	Command    string          `json:"command"`
	Data       json.RawMessage `json:"data,omitempty"`
	Identifier Identifier      `json:"identifier"`
	ErrorChan  chan error      `json:"-"`
}

// Identifier is the channel name itself
type Identifier struct {
	Channel string
}

type innerIdentifier struct {
	Channel string `json:"channel"`
}

// NewSubscription sugar wrapper for subscription command
func NewSubscription(chanName string) *Command {
	return NewCommand("subscribe", chanName)
}

// CancelSubscription sugar wrapper for subscription command
func CancelSubscription(chanName string) *Command {
	return NewCommand("unsubscribe", chanName)
}

// NewCommand instantiates a new command
func NewCommand(command, chanName string) *Command {
	return &Command{
		Command: command,
		Identifier: Identifier{
			Channel: chanName,
		},
	}
}

// MarshalJSON encodes the Identifier from JSON
func (i *Identifier) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(innerIdentifier{
		Channel: i.Channel,
	})
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(b))
}

// UnmarshalJSON decodes the Identifier from JSON. Because the inner identifier is double encoded
func (i *Identifier) UnmarshalJSON(data []byte) error {
	str := ""
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	inner := innerIdentifier{}
	if err := json.Unmarshal([]byte(str), &inner); err != nil {
		return err
	}
	i.Channel = inner.Channel
	return nil
}
