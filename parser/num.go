package main

import (
	"fmt"
	"strconv"
	"unicode"
)

type Num struct{ val int }

func NewNum(s string) (*Num, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, fmt.Errorf("invalid num: %s", s)
	}
	return &Num{n}, nil
}

func (n *Num) Type() string   { return "num" }
func (n *Num) Value() any     { return n.val }
func (n *Num) String() string { return fmt.Sprintf("num:%d", n.val) }

func (n *Num) Check(r rune) bool { return unicode.IsDigit(r) }
func (n *Num) Parse(s string) (Data, string, error) {
	if err := checkInit(n, s); err != nil {
		return handleError(err)
	}

	var res string
	for i, r := range s {
		switch {
		case n.Check(r):
			res += string(r)
			continue
		case unicode.IsSpace(r):
			if res == "" {
				return handleError(fmt.Errorf("empty num"))
			}
			num, err := NewNum(res)
			if err != nil {
				return handleError(fmt.Errorf("invalid num: %s", res))
			}
			return num, s[i:], nil
		default:
			if len(res) > 0 {
				num, err := NewNum(res)
				if err != nil {
					return handleError(err)
				}
				return num, s[i:], nil
			}
			return handleError(NewUnexpectedCharErr("num:default", r))
		}
	}

	num, err := strconv.Atoi(res)
	if err != nil {
		return handleError(fmt.Errorf("invalid num: %s", res))
	}
	return &Num{num}, "", nil
}
