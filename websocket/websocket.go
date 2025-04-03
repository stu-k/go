package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var ErrCloseExpected = fmt.Errorf("close expected")
var ErrCloseUnxpected = fmt.Errorf("close unexpected")

type Conn interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}

func Dial(url string, h http.Header) (Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(url, h)
	return &gorillaConn{conn}, err
}

func Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	conn, err := upgrader.Upgrade(w, r, h)
	if err != nil {
		return nil, err
	}
	return &gorillaConn{conn}, nil
}

type gorillaConn struct {
	conn *websocket.Conn
}

func (c *gorillaConn) ReadJSON(v any) error {
	err := c.conn.ReadJSON(v)
	if websocket.IsCloseError(err) {
		return ErrCloseExpected
	}
	if websocket.IsUnexpectedCloseError(err) {
		return ErrCloseUnxpected
	}
	return err
}

func (c *gorillaConn) WriteJSON(v any) error {
	err := c.conn.WriteJSON(v)
	if websocket.IsCloseError(err) {
		return ErrCloseExpected
	}
	if websocket.IsUnexpectedCloseError(err) {
		return ErrCloseUnxpected
	}
	return err
}

func (c *gorillaConn) Close() error {
	return c.conn.Close()
}
