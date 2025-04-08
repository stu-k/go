package syntax

import (
	"fmt"
	"strconv"
	"strings"
)

type seqStrArgs struct {
	rule        *Rule
	char        rune
	count       int
	cap, usecap bool
}

func NewSequenceFromStrs(name string, parts ...string) (*Sequence, error) {
	return newSequenceFromStrs(name, defaultRulemap, parts...)
}

func newSequenceFromStrs(name string, pmap map[string]*Rule, parts ...string) (*Sequence, error) {
	errFn := func(arg string, i, j int) error {
		return fmt.Errorf(
			"error creating sequence from strs: invalid arg \"%v\" in segment %d, arg %d",
			arg, i, j,
		)
	}

	sq := NewSequence(name)

	if len(parts) == 0 {
		return nil, fmt.Errorf("error creating sequence from strs: invalid string \"%s\"", parts)
	}

	for i, part := range parts {

		var sqa seqStrArgs

		args := strings.Split(part, ", ")
		for j, arg := range args {
			if len(arg) == 0 || len(arg) == 1 {
				return nil, errFn(arg, i, j)
			}

			r := rune(arg[1])
			switch rune(arg[0]) {

			case 'r':
				seq, ok := pmap[arg[1:]]
				if !ok || seq == nil {
					return nil, errFn(arg, i, j)
				}
				sqa.rule = seq
				continue

			case '#':
				ct, err := strconv.Atoi(arg[1:])
				if err != nil {
					return nil, errFn(arg, i, j)
				}
				sqa.count = ct
				continue

			case 'c':
				sqa.char = r
				continue

			case 'g':
				if r == '0' {
					sqa.cap = false
					sqa.usecap = true
					continue
				}
				if r == '1' {
					sqa.cap = true
					sqa.usecap = true
					continue
				}
				return nil, errFn(arg, i, j)

			default:
				return nil, errFn(arg, i, j)

			}
		}

		rule := NewRule()
		if sqa.rule != nil {
			rule = sqa.rule
		}
		if sqa.count != 0 {
			rule = rule.Count(sqa.count)
		}
		if sqa.char != 0 {
			rule = rule.Check(func(r rune) bool { return r == sqa.char })
		}
		if sqa.usecap {
			rule = rule.Capture(sqa.cap)
		}
		rule = rule.Named(part)
		if rule.IsAny() {
			return nil, fmt.Errorf("error creating sequence from strs: can't add empty rule \"%s\"", part)
		}
		sq.Add(rule)
	}

	if sq.Len() == 0 {
		return nil, fmt.Errorf("error creating sequence from strs: can't use empty sequence")
	}

	return sq, nil
}
