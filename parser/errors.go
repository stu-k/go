package main

import "fmt"

type ParseError struct{ msg string }

func (p ParseError) Error() string {
	return p.msg
}
func NewUnexpectedCharErr(which string, char rune) ParseError {
	return ParseError{fmt.Sprintf("unexpected char in [%s]: \"%s\"", which, string(char))}
}

func NewExpectationErr(got rune, want rune) ParseError {
	return ParseError{fmt.Sprintf("expected \"%s\"; got \"%s\"", string(got), string(want))}
}

func NewSingleExpectationErr(want rune) ParseError {
	return ParseError{fmt.Sprintf("expected \"%s\"", string(want))}
}

func handleError(err error) (Data, string, error) {
	return nil, "", err
}

func checkInit(data interface {
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
