package bus

import (
	"errors"
	"fmt"
	"testing"
)

type TestQuery struct {
	ID       int64
	Response string
}

func TestQueryHandlerReturnsError(t *testing.T) {
	bus := NewBus()

	bus.AddRequestHandler(func(q *TestQuery) error {
		return errors.New("handler error")
	})

	err := bus.SendRequest(&TestQuery{})

	if err == nil {
		t.Fatal("Send query failed " + err.Error())
	} else {
		t.Log("Handler error received ok")
	}
}

func TestQueryHandlerReturn(t *testing.T) {
	bus := NewBus()

	bus.AddRequestHandler(func(q *TestQuery) error {
		q.Response = "hello from handler"
		return nil
	})

	q := &TestQuery{}
	err := bus.SendRequest(q)

	if err != nil {
		t.Fatal("Send query failed " + err.Error())
	} else if q.Response != "hello from handler" {
		t.Fatal("Failed to get response from handler")
	}
}

func TestEventListeners(t *testing.T) {
	bus := NewBus()
	count := 0

	bus.AddEventListener(func(query *TestQuery) error {
		count++
		return nil
	})

	bus.AddEventListener(func(query *TestQuery) error {
		count += 10
		return nil
	})

	err := bus.BroadcastEvent(&TestQuery{})

	if err != nil {
		t.Fatal("Publish event failed " + err.Error())
	} else if count != 11 {
		t.Fatal(fmt.Sprintf("Publish event failed, listeners called: %v, expected: %v", count, 11))
	}
}
