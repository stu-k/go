package main

import (
	"fmt"
	"unicode"

	"github.com/stu-k/go/parser/errors"
)

type Obj struct{ val map[string]Data }

func NewObject(val map[string]Data) *Obj {
	return &Obj{val}
}

func (o *Obj) Type() string { return "object" }
func (o *Obj) Value() any   { return o.val }
func (o *Obj) String() string {
	sofar := "obj:{"
	for k, v := range o.val {
		sofar += fmt.Sprintf(" %s: %v", k, v)
	}
	return sofar + " }"
}

func (o *Obj) Check(r rune) bool { return r == '{' }
func (o *Obj) Parse(s string) (Data, string, error) {
	if err := errors.CheckInit(o, s); err != nil {
		return errors.HandeleError(err)
	}

	toparse := s[1:]

	res := make(map[string]Data)

	var key string
	isColon := false
	isComma := false
	for i := 0; i < len(toparse); i++ {
		r := rune(toparse[i])
		switch {
		case unicode.IsSpace(r):
			continue
		case r == '}':
			// can't end on key or comma
			if key != "" || isComma {
				return errors.HandeleError(errors.NewUnexpectedCharErr('}'))
			}
			return NewObject(res), toparse[i+1:], nil
		case r == ':':
			// can't use colon without key
			if key == "" {
				return errors.HandeleError(errors.NewUnexpectedCharErr(':'))
			}
			isColon = true
			continue
		case r == ',':
			// comma comes after a complete kv set
			if len(res) == 0 || key != "" {
				return errors.HandeleError(errors.NewUnexpectedCharErr(','))
			}
			isComma = true
			continue
		default:
			// need colon after key
			if key != "" && !isColon {
				return errors.HandeleError(errors.NewExpectedCharErr(':'))
			}

			// parse token
			data, rest, err := parse(toparse[i:], mainOpts)
			if err != nil {
				return errors.HandeleError(err)
			}

			// set value
			if key != "" {
				// need comma between kv sets
				if len(res) > 0 && !isComma {
					return errors.HandeleError(errors.NewExpectedCharErr(','))
				}
				res[key] = data
				key = ""
				isColon = false
				isComma = false
				toparse = rest
				i = -1
				continue
			}

			// set key
			switch data.Type() {
			case "str":
				k, err := handleSetKey(data)
				if err != nil {
					return errors.HandeleError(err)
				}
				key = k
				toparse = rest
				i = -1
				continue
			case "var":
				k, err := handleSetKey(data)
				if err != nil {
					return errors.HandeleError(err)
				}
				key = k
				toparse = rest
				i = -1
				continue
			default:
				return errors.HandeleError(fmt.Errorf("invalid obj key type: %s", data.Type()))
			}
		}
	}

	return errors.HandeleError(errors.NewExpectedCharErr('}'))
}

func handleSetKey(data Data) (string, error) {
	k, ok := data.Value().(string)
	if !ok {
		// value of data isn't a string, internal error
		// type "str".val != string
		// type "token".val != string
		panic(fmt.Errorf("invalid obj key type: %v", data.Value()))
	}
	if k == "" {
		return "", fmt.Errorf("obj key can not be empty")
	}
	return k, nil
}
