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

func (n Num) Check(r rune) bool { return unicode.IsDigit(r) }
func (n Num) Parse(s string) (Data, string, error) {
	if err := checkInit(n, s); err != nil {
		panic(err)
	}

	var sofar string
	for i, r := range s {
		switch {
		case n.Check(r):
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
			return num, s[i:], nil
		default:
			if len(sofar) > 0 {
				num, err := NewNum(sofar)
				if err != nil {
					return handleError(err)
				}
				return num, s[i:], nil
			}
			return handleError(NewUnexpectedTokenErr("num:default", r))
		}
	}

	num, err := strconv.Atoi(sofar)
	if err != nil {
		return handleError(fmt.Errorf("invalid num: %s", sofar))
	}
	return Num{num}, "", nil
}
