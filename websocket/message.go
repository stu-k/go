package websocket

import (
	"encoding/json"
	"fmt"
)

// Message represents a structured message between client and server
type Message struct {
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
	Meta  json.RawMessage `json:"meta"`
	Error string          `jons:"error"`
}

func NewMessage(t string, d interface{}, meta json.RawMessage, e error) (*Message, error) {
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
		Meta:  meta,
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
