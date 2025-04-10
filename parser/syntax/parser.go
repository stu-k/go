package syntax

import "fmt"

type Parsable interface {
	Parse(string) (*ParseResult, error)
	Name() string
}

type Parser struct {
	name  string
	seqs  map[string]*Sequence
	rules map[string]*Rule
}

func NewParser(name string) *Parser {
	return &Parser{
		name: name,
		seqs: make(map[string]*Sequence),
	}
}

func (p *Parser) Name() string {
	return p.name
}

func (p *Parser) NewSeq(n string, s ...Sequencer) (*Sequence, error) {
	_, ok := p.seqs[n]
	if ok {
		return nil, fmt.Errorf("[parser:newseq] sequence with name \"%v\" already exists", n)
	}

	ps := make([]Parsable, 0)
	for _, sq := range s {
		ps = append(ps, sq.Seq())
	}
	seq := NewSequence(n, ps...)
	p.seqs[n] = seq
	return seq, nil
}

func (p *Parser) NewRule(n string) (*Rule, error) {
	_, ok := p.rules[n]
	if ok {
		return nil, fmt.Errorf("[parser:newrule] rule with name \"%v\" already exists", n)
	}
	r := NewRule(n)
	p.rules[n] = r
	return r, nil
}
