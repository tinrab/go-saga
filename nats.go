package saga

import (
	"context"
	"encoding/json"
	"github.com/nats-io/go-nats"
)

type natsBus struct {
	conn *nats.Conn
}

func NewNatsBus(conn *nats.Conn) *natsBus {
	return &natsBus{
		conn: conn,
	}
}

func (n *natsBus) Close() {
	n.conn.Close()
}

func (n *natsBus) PublishEvent(event string, data interface{}) error {
	bd, err := n.encode(data)
	if err != nil {
		return err
	}
	return n.conn.Publish(event, bd)
}

func (n *natsBus) ReceiveEvent(ctx context.Context, event string, res interface{}) error {
	sub, err := n.conn.SubscribeSync(event)
	if err != nil {
		return err
	}
	msg, err := sub.NextMsgWithContext(ctx)
	if err != nil {
		return err
	}

	if err = n.decode(msg.Data, &res); err != nil {
		return err
	}

	return nil
}

func (n *natsBus) HandleEvent(event string, handler EventHandler) error {
	_, err := n.conn.Subscribe(event, func(msg *nats.Msg) {
		handler(msg.Subject, msg.Data)
	})
	return err
}

func (n *natsBus) Request(ctx context.Context, event string, data interface{}, res interface{}) error {
	bd, err := n.encode(data)
	if err != nil {
		return err
	}

	msg, err := n.conn.RequestWithContext(ctx, event, bd)
	if err != nil {
		return err
	}

	if err := n.decode(msg.Data, &res); err != nil {
		return err
	}

	return nil
}

func (n *natsBus) Reply(event string, handler RequestHandler) error {
	_, err := n.conn.Subscribe(event, func(msg *nats.Msg) {
		res := handler(event, msg.Data)
		rb, err := n.encode(res)
		if err != nil {
			panic(err)
		}

		if err = n.conn.Publish(msg.Reply, rb); err != nil {
			panic(err)
		}
	})
	return err
}

func (n *natsBus) encode(v interface{}) ([]byte, error) {
	if enc, ok := v.(Encoder); ok {
		return enc.Encode()
	}

	return json.Marshal(v)
}

func (n *natsBus) decode(data []byte, v interface{}) error {
	if dec, ok := v.(Decoder); ok {
		return dec.Decode(data)
	}

	return json.Unmarshal(data, v)
}
