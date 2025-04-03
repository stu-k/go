package websocket

import (
	"encoding/json"
	"fmt"
)

// Message represents a structured message between client and server
type Message struct {
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
	Error string          `jons:"error"`
}

func NewMessage(t string, d interface{}, e error) (*Message, error) {
	data, err := json.Marshal(d)
	if err != nil {
		return nil, err
	}

	errS := ""
	if e != nil {
		errS = e.Error()
	}

	return &Message{
		Type:  t,
		Data:  data,
		Error: errS,
	}, nil
}

func (m *Message) String() string {
	s := struct {
		Type  string
		Data  string
		Error string
	}{
		m.Type,
		string(m.Data),
		m.Error,
	}
	return fmt.Sprintf("%+v", s)
}
