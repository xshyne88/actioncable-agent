package actioncable

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestMarshalIdentifier(t *testing.T) {
	chanName := "AgentChannel"
	i := &Identifier{
		Channel: chanName,
	}
	exp := []byte(fmt.Sprintf(`"{\"channel\":\"%s\"}"`, chanName))

	act, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}

	assert(reflect.DeepEqual(act, exp), true, t)
}

func TestUnmarshalIdentifier(t *testing.T) {
	r := &Identifier{
		Channel: "AgentChannel",
	}
	bytes, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	var exp Identifier
	if err := json.Unmarshal(bytes, &exp); err != nil {
		t.Fatal(err)
	}

	assertDeep(*r, exp, t)
}

func TestNewCommand(t *testing.T) {
	chanName := "AgentChannel"
	exp := &Command{
		Command: "subscribe",
		Identifier: Identifier{
			Channel: chanName,
		},
	}

	act := NewCommand("subscribe", chanName)

	assertDeep(act, exp, t)
}

func TestNewSubscriptionCommand(t *testing.T) {
	chanName := "AgentChannel"
	exp := &Command{
		Command: "subscribe",
		Identifier: Identifier{
			Channel: chanName,
		},
	}

	act := NewSubscription(chanName)

	assertDeep(act, exp, t)
}

func TestCancelSubscriptionCommand(t *testing.T) {
	chanName := "AgentChannel"
	exp := &Command{
		Command: "unsubscribe",
		Identifier: Identifier{
			Channel: chanName,
		},
	}

	act := CancelSubscription(chanName)

	assertDeep(act, exp, t)
}
