package saga

import (
	"context"
)

type EventHandler func(event string, data []byte)

type RequestHandler func(event string, data []byte) interface{}

type Bus interface {
	PublishEvent(event string, data interface{}) error
	ReceiveEvent(ctx context.Context, event string, res interface{}) error
	HandleEvent(event string, handler EventHandler)
	Request(ctx context.Context, event string, data interface{}, res interface{}) error
	Reply(event string, handler RequestHandler)
}
