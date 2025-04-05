package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type Num struct{ val int }

func NewNum(s string) (Num, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return Num{}, fmt.Errorf("invalid num: %s", s)
	}
	return Num{n}, nil
}

func (n Num) Type() string   { return "num" }
func (n Num) Value() any     { return n.val }
func (n Num) String() string { return fmt.Sprintf("num:%d", n.val) }

func isNum(r rune) bool { return unicode.IsDigit(r) }
func parseNum(input string) (Data, string, error) {
	var sofar string
	for i, r := range input {
		switch {
		case isNum(r):
			sofar += string(r)
			continue
		case unicode.IsSpace(r):
			if sofar == "" {
				return handleError(fmt.Errorf("empty num"))
			}
			num, err := NewNum(sofar)
			if err != nil {
				return handleError(fmt.Errorf("invalid num: %s", sofar))
			}
			return num, input[i:], nil
		default:
			if len(sofar) > 0 {
				num, err := NewNum(sofar)
				if err != nil {
					return handleError(err)
				}
				return num, input[i:], nil
			}
			return handleError(fmt.Errorf("invalid char in num: %s", string(r)))
		}
	}
	n, err := strconv.Atoi(sofar)
	if err != nil {
		return handleError(fmt.Errorf("invalid num: %s", sofar))
	}
	return Num{n}, "", nil
}
