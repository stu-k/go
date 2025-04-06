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

func HandeleError(err error) (data, string, error) {
	return nil, "", err
}

func CheckInit(data interface {
	Check(rune) bool
	Type() string
}, s string) error {
	if s == "" {
		return fmt.Errorf("%s init with \"\"", data.Type())
	}
	if !data.Check(rune(s[0])) {
		return fmt.Errorf("%s init with \"%s\"", data.Type(), string(s[0]))
	}
	return nil
}
