package main

import (
	"fmt"
	"strings"
)

type anyfn func(...Data) (Data, error)

type Op struct {
	val string
	fn  anyfn
	air int
}

func NewOp(v string, fn anyfn, air int) Op { return Op{v, fn, air} }
func (o Op) Type() string                  { return "op" }
func (o Op) Value() any                    { return o.val }
func (o Op) String() string                { return fmt.Sprintf("op:%s", o.val) }
func (o Op) Exec(d ...Data) (Data, error)  { return o.fn(d...) }

func isOp(s string) bool {
	for k, _ := range opMap {
		if strings.HasPrefix(k, s) {
			return true
		}
	}
	return false
}

var opMap = map[string]anyfn{
	"*m*": opMult,
}

func opMult(d ...Data) (Data, error) { return nil, nil }

func parseOp(s string) (Data, string, error) {
	if s == "" {
		panic(fmt.Errorf("op init with \"\""))
	} else if !isOp(string(s[0])) {
		panic(fmt.Errorf("op init with \"%s\"", string(s[0])))
	}

	// sofar := string(input[0])
	for i := 0; i < len(s); i++ {

	}

	return nil, "", nil
}
