package parse

import (
	"fmt"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

type Arr struct{ val []Data }

func (a *Arr) Type() string { return "arr" }
func (a *Arr) Value() any   { return a.val }
func (a *Arr) String() string {
	sofar := "arr:["
	for _, d := range a.val {
		sofar += fmt.Sprintf(" %v", d)
	}
	sofar += " ]"
	return sofar
}

func (a *Arr) Check(r rune) bool { return r == '[' }
func (a *Arr) Parse(s string) (Data, string, error) {
	if err := errors.CheckInit(a, s); err != nil {
		return errors.HandeleError(err)
	}

	toparse := s[1:]
	res := make([]Data, 0)

	lastWasComma := false
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == ']':
			if lastWasComma {
				return errors.HandeleError(errors.NewExpectedCharErr(']'))
			}
			return &Arr{res}, toparse[i+1:], nil
		case r == ',':
			if len(res) == 0 {
				return errors.HandeleError(errors.NewUnexpectedCharErr(r))
			}
			lastWasComma = true
			continue
		default:
			if len(res) > 0 && !lastWasComma {
				return errors.HandeleError(errors.NewExpectedCharErr(']'))
			}

			data, rest, err := parse(toparse[i:], mainOpts)
			if err != nil {
				return errors.HandeleError(err)
			}

			res = append(res, data)
			if rest != "" {
				toparse = rest
				i = -1
			}
			lastWasComma = false
			continue
		}
	}
	return errors.HandeleError(errors.NewExpectedCharErr(']'))
}
