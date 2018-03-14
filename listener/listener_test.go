package listener

import (
	"encoding/json"
	"testing"
)

type TestMessage struct {
	S string `json:"s"`
	I int    `json:"i"`
}

func TestQueryHandlerReturn(t *testing.T) {
	bus := NewListener()
	var s string
	var i int

	bus.AddHandler(func(q *TestMessage) error {
		s = q.S
		i = q.I
		return nil
	})

	q := &TestMessage{S: "test", I: 1}
	b, _ := json.Marshal(q)
	err := bus.Handle("TestMessage", b)

	if err != nil {
		t.Fatal("Handle failed " + err.Error())
	} else if s != "test" || i != 1 {
		t.Fatalf("Failed to get response from handler %v, %v", s, i)
	}
}
