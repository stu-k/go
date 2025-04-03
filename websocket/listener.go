package websocket

import (
	"context"
	"fmt"
)

type listener interface {
	Message() chan *Message
	Error() <-chan error
	Listen()
}

type inListener struct {
	message chan *Message
	err     chan error
	conn    Conn
	ctx     context.Context
}

func NewInListener(ctx context.Context, conn Conn) (listener, error) {
	if conn == nil {
		return nil, fmt.Errorf("error creating in listener: conn is nil")
	}

	in := &inListener{
		message: make(chan *Message),
		err:     make(chan error),
		conn:    conn,
		ctx:     ctx,
	}

	return in, nil
}
func (l *inListener) Message() chan *Message {
	return l.message
}

func (l *inListener) Error() <-chan error {
	return l.err
}

func (l *inListener) Listen() {
	go func(c chan<- *Message, e chan<- error) {
		for {
			select {
			case <-l.ctx.Done():
				close(c)
				close(e)
				return
			default:
				var msg *Message
				err := l.conn.ReadJSON(&msg)
				if err != nil {
					e <- err
					continue
				}
				c <- msg
				select {
				case <-l.ctx.Done():
					close(c)
					close(e)
					return
				default:
					continue
				}
			}
		}
	}(l.message, l.err)
}

type outListener struct {
	message chan *Message
	err     chan error
	conn    Conn
	ctx     context.Context
}

func NewOutListener(ctx context.Context, conn Conn) (listener, error) {
	if conn == nil {
		return nil, fmt.Errorf("error creating out listener: conn is nil")
	}

	out := &outListener{
		message: make(chan *Message),
		err:     make(chan error),
		conn:    conn,
		ctx:     ctx,
	}

	return out, nil
}
func (l *outListener) Message() chan *Message {
	return l.message
}

func (l *outListener) Error() <-chan error {
	return l.err
}

func (l *outListener) Listen() {
	go func(c <-chan *Message, e chan<- error) {
		for {
			select {
			case <-l.ctx.Done():
				close(e)
				return
			case msg := <-c:
				err := l.conn.WriteJSON(&msg)
				if err != nil {
					e <- err
					continue
				}
			}
		}
	}(l.message, l.err)
}
