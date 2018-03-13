package bus

import (
	"fmt"
	"reflect"
)

// Request is a generic request
type Request interface{}

// Response is a response to a request
type Response interface{}

// Event is a generic event
type Event interface{}

// HandlerFunc is a handler
type HandlerFunc interface{}

// Bus defines an interface for a bus
type Bus interface {
	SendRequest(r Request) (Response, error)
	BroadcastEvent(e Event) error

	AddRequestHandler(handler HandlerFunc)
	AddEventListener(handler HandlerFunc)
}

type busImpl struct {
	handlers  map[string]HandlerFunc
	listeners map[string][]HandlerFunc
}

// NewBus creates a Bus instance
func NewBus() Bus {
	return &busImpl{
		handlers:  make(map[string]HandlerFunc),
		listeners: make(map[string][]HandlerFunc),
	}
}

func (b *busImpl) SendRequest(r Request) (Response, error) {
	var rType = reflect.TypeOf(r).Elem().Name()

	var handler = b.handlers[rType]
	if handler == nil {
		return nil, fmt.Errorf("handler not found for %s", rType)
	}

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(r)

	ret := reflect.ValueOf(handler).Call(params)
	rsp := ret[0].Interface()
	err := ret[1].Interface()
	if err != nil {
		return nil, err.(error)
	}

	return rsp.(Response), nil
}

func (b *busImpl) BroadcastEvent(e Event) error {
	var eType = reflect.TypeOf(e).Elem().Name()
	var listeners = b.listeners[eType]

	var params = make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(e)

	for _, listenerHandler := range listeners {
		ret := reflect.ValueOf(listenerHandler).Call(params)
		err := ret[0].Interface()
		if err != nil {
			return err.(error)
		}
	}

	return nil
}

func (b *busImpl) AddRequestHandler(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	queryTypeName := handlerType.In(0).Elem().Name()
	b.handlers[queryTypeName] = handler
}

func (b *busImpl) AddEventListener(handler HandlerFunc) {
	handlerType := reflect.TypeOf(handler)
	eventName := handlerType.In(0).Elem().Name()
	_, exists := b.listeners[eventName]
	if !exists {
		b.listeners[eventName] = make([]HandlerFunc, 0)
	}
	b.listeners[eventName] = append(b.listeners[eventName], handler)
}
