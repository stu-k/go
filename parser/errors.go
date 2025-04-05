package main

import "fmt"

type ParseError struct{ msg string }

func (p ParseError) Error() string {
	return p.msg
}
func NewUnexpectedTokenErr(which string, token rune) ParseError {
	return ParseError{fmt.Sprintf("unexpected token in [%s]: %s", which, string(token))}
}

func NewExpectationErr(got rune, want rune) ParseError {
	return ParseError{fmt.Sprintf("expected \"%s\"; got \"%s\"", string(got), string(want))}
}

func NewSingleExpectationErr(want rune) ParseError {
	return ParseError{fmt.Sprintf("expected \"%s\"", string(want))}
}

func handleError(err error) (Data, string, error) {
	return DataUnknown{}, "", err
}
