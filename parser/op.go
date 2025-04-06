package main

import (
	"fmt"

	"github.com/stu-k/go/parser/errors"
)

type anyfn func(...Data) (Data, error)

type Op struct {
	val  string
	fn   anyfn
	args []Data
	lair uint
	rair uint
}

func NewOp(val string, fn anyfn, args []Data, air uint) *Op {
	return &Op{val, fn, args, 0, air}
}
func (o *Op) Type() string   { return "op" }
func (o *Op) Value() any     { return o.val }
func (o *Op) String() string { return fmt.Sprintf("op:%s[%d]", o.val, o.rair) }

var noop = func(d ...Data) (Data, error) {
	return nil, nil
}

var opBang = NewOp("!", noop, []Data{}, 2)

var opMap = map[string]*Op{
	// "~":  *Op,
	"!": opBang,
	// "@":  *Op,
	// "#":  *Op,
	// "$":  *Op,
	// "%":  *Op,
	// "^":  *Op,
	// "&":  *Op,
	// "*":  *Op,
	// "-":  *Op,
	// "+":  *Op,
	// "=":  *Op,
	// "|":  *Op,
	// "\\": *Op,
	// "<":  *Op,
	// ">":  *Op,
	// ".":  *Op,
	// "?":  *Op,
	// "/":  *Op,
}

func (op Op) Check(r rune) bool {
	c := string(r)
	for k, _ := range opMap {
		if string(k[0]) == c {
			return true
		}
	}
	return false
}
func (op *Op) Parse(s string) (Data, string, error) {
	if err := errors.CheckInit(op, s); err != nil {
		return errors.HandelError(err)
	}

	// var res string
	// for i, char := range s {
	// 	c := string(char)
	// 	for k, _ := range opMap {
	// 		if !strings.HasPrefix(k, res+c) {
	// 			fn, ok := opMap[res]
	// 			if !ok {
	// 				panic(fmt.Errorf("op \"res\" not found"))
	// 			}
	// 		}
	// 	}
	// }
	return nil, "", nil
	// return NewOp(string(s[0]), fn, []Data{}, 2), s[1:], nil
}
