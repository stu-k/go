package errors

import "fmt"

var ErrUnexpectedChar = fmt.Errorf("unexpected char")

type UnexpectedCharErr struct {
	char rune
}

func (e UnexpectedCharErr) Error() string {
	return fmt.Sprintf("unexpected char \"%s\"", string(e.char))
}
func (e UnexpectedCharErr) Unwrap() error { return ErrUnexpectedChar }
func NewUnexpectedCharErr(char rune) error {
	return UnexpectedCharErr{char}
}

var ErrExpectedChar = fmt.Errorf("expected vhar")

type ExpectedCharErr struct {
	char rune
}

func NewExpectedCharErr(char rune) error {
	return ExpectedCharErr{char}
}

func (e ExpectedCharErr) Error() string {
	return fmt.Sprintf("expected \"%s\"", string(e.char))
}

func (e ExpectedCharErr) Unwrap() error {
	return ErrExpectedChar
}

var ErrEndOfInput = fmt.Errorf("end of input")

type EndOfInputErr struct{}

func (e EndOfInputErr) Error() string {
	return ErrEndOfInput.Error()
}
func (e EndOfInputErr) Unwrap() error { return ErrEndOfInput }
func NewEndOfInputErr() error {
	return EndOfInputErr{}
}

type data interface {
	Type() string
	Value() any
	String() string
}

var ErrBadMatch = fmt.Errorf("bad match")

type BadMatchErr struct {
	name string
	bad  string
	tag  string
}

func (e BadMatchErr) Error() string {
	info := e.tag
	if info != "" {
		info = "[" + info + "]"
	}
	return fmt.Sprintf("bad match for %s%v: \"%s\"", e.name, info, e.bad)
}
func (e BadMatchErr) Unwrap() error { return ErrBadMatch }
func NewBadMatchErr(name, bad string, v ...string) error {
	tag := ""
	if len(v) > 0 {
		tag = v[0]
	}
	return BadMatchErr{name, bad, tag}
}

func HandleError(err error) (data, string, error) {
	return nil, "", err
}

func CheckInit(name string, s string, check func(rune) bool) error {
	if s == "" {
		return fmt.Errorf("%s init with \"\"", name)
	}
	if !check(rune(s[0])) {
		return fmt.Errorf("%s init with \"%s\"", name, string(s[0]))
	}
	return nil
}
