package parse

import (
	"fmt"

	"github.com/stu-k/go/parser/errors"
)

type Str struct{ val string }

func NewStr(s string) *Str    { return &Str{s} }
func (s *Str) Type() string   { return "str" }
func (s *Str) Value() any     { return s.val }
func (s *Str) String() string { return fmt.Sprintf("str:\"%s\"", s.val) }

func (s *Str) Check(r rune) bool { return r == '"' }
func (str *Str) Parse(s string) (Data, string, error) {
	if err := errors.CheckInit(str, s); err != nil {
		return errors.HandeleError(err)
	}

	toparse := s[1:]

	var res string
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case r == '"':
			return NewStr(res), toparse[i+1:], nil
		default:
			res += string(r)
			continue
		}
	}

	return errors.HandeleError(errors.NewExpectedCharErr('"'))
}
