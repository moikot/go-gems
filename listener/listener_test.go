package listener

import (
	"encoding/json"
	"testing"
)

type TestMessage struct {
	s string
	i int
}

func TestQueryHandlerReturn(t *testing.T) {
	bus := NewListener()
	var s string
	var i int

	bus.AddHandler(func(q *TestMessage) error {
		s = q.s
		i = q.i
		return nil
	})

	q := &TestMessage{s: "test", i: 1}
	b, _ := json.Marshal(q)

	err := bus.Handle("TestMessage", b)

	if err != nil {
		t.Fatal("Handle failed " + err.Error())
	} else if s != "test" || i != 1 {
		t.Fatal("Failed to get response from handler")
	}
}
