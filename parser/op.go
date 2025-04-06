package main

import (
	"fmt"
)

type anyfn func(...Data) (Data, error)

type Op struct {
	val  string
	fn   anyfn
	args []Data
	air  uint
}

func NewOp(val string, fn anyfn, args []Data, air uint) Op {
	return Op{val, fn, args, air}
}
func (o Op) Type() string   { return "op" }
func (o Op) Value() any     { return o.val }
func (o Op) String() string { return fmt.Sprintf("op:\"%s[%d]\"", o.val, o.air) }

var noop = func(d ...Data) (Data, error) {
	return nil, nil
}

var opMap = map[rune]anyfn{
	'~':  noop,
	'!':  noop,
	'@':  noop,
	'#':  noop,
	'$':  noop,
	'%':  noop,
	'^':  noop,
	'&':  noop,
	'*':  noop,
	'-':  noop,
	'+':  noop,
	'=':  noop,
	'|':  noop,
	'\\': noop,
	'<':  noop,
	'>':  noop,
	'.':  noop,
	'?':  noop,
	'/':  noop,
}

func (s Op) Check(r rune) bool {
	_, ok := opMap[r]
	return ok
}
func (op Op) Parse(s string) (Data, string, error) {
	if err := checkInit(op, s); err != nil {
		panic(err)
	}

	fn, ok := opMap[rune(s[0])]
	if !ok {
		panic(fmt.Errorf("no op found for rune %s", string(s[0])))
	}

	return NewOp(string(s[0]), fn, []Data{}, 2), s[1:], nil
}
