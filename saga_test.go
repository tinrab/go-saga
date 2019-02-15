package saga

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/go-nats"
	"testing"
)

type vector struct {
	X int
	Y int
}

type result struct {
	Value int
}

func (v *vector) Decode(data []byte) error {
	return json.Unmarshal(data, v)
}

func TestSanity(t *testing.T) {
	conn, err := nats.Connect("localhost:4222")
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	nb := NewNatsBus(conn)

	if err := nb.Reply("add", func(event string, data []byte) interface{} {
		p := &vector{}
		if err := p.Decode(data); err != nil {
			panic(err)
		}

		return result{p.X + p.Y}
	}); err != nil {
		t.Fatal(err)
	}

	var res result

	if err := nb.Request(context.TODO(), "add", &vector{3, 2}, &res); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println("res:", res.Value)
	}
}
