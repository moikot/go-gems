package listener

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Request is a generic request
type Request interface{}

// Event is a generic event
type Event interface{}

// HandlerFunc is a handler
type HandlerFunc interface{}

// Bus defines an interface for a bus
type Listener interface {
	Handle(messageType string, data []byte) error

	AddHandler(handler HandlerFunc)
}

type listener struct {
	handlers map[string]HandlerFunc
}

// NewBus creates a Bus instance
func NewListener() Listener {
	return &listener{
		handlers: make(map[string]HandlerFunc),
	}
}

func (b *listener) Handle(messageType string, data []byte) error {
	var rType = messageType

	var handler = b.handlers[rType]
	if handler == nil {
		return fmt.Errorf("handler not found for %s", rType)
	}

	// Deserialize
	handlerType := reflect.TypeOf(handler)
	elem := handlerType.In(0).Elem()
	v := reflect.New(elem)

	if err := json.Unmarshal(data, v.Interface()); err != nil {
		return fmt.Errorf("unable to marshal %s, %v", elem.Name(), err)
	}

	var params = make([]reflect.Value, 1)
	params[0] = v

	ret := reflect.ValueOf(handler).Call(params)
	err := ret[0].Interface()
	if err == nil {
		return nil
	} else {
		return err.(error)
	}
}

func (b *listener) AddHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(0).Elem().Name()
	b.handlers[queryTypeName] = handler
}
