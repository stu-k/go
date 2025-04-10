package syntax

import "fmt"

type Parsable interface {
	Parse(string) (*Result, error)
	Name() string
}

type Parser struct {
	name string
	seqs map[string]Parsable
}

func NewParser(name string) *Parser {
	return &Parser{
		name: name,
		seqs: make(map[string]Parsable),
	}
}

func (p *Parser) Name() string {
	return p.name
}

func (p *Parser) NewSeq(n string, r ...any) error {
	new, err := p.newSeq(n, r...)
	if err != nil {
		return err
	}
	_, ok := p.seqs[n]
	if ok {
		fmt.Printf("[parser:newseq] overwriting sequence with name \"%v\"", n)
	}
	p.seqs[n] = new
	return nil
}

func (p *Parser) NewPickOne(n string, r ...any) error {
	new, err := p.newSeq(n, r...)
	if err != nil {
		return err
	}
	_, ok := p.seqs[n]
	if ok {
		fmt.Printf("[parser:pickone] overwriting sequence with name \"%v\"", n)
	}
	p.seqs[n] = new.PickOne()
	return nil
}

func (p *Parser) NewUntilFail(n string, r ...any) error {
	new, err := p.newSeq(n, r...)
	if err != nil {
		return err
	}
	_, ok := p.seqs[n]
	if ok {
		fmt.Printf("[parser:untilfail] overwriting sequence with name \"%v\"", n)
	}
	p.seqs[n] = new.UntilFail()
	return nil
}
func (p *Parser) newSeq(n string, r ...any) (*Sequence, error) {
	res := make([]Parsable, 0)
	for i, v := range r {
		par, ok := v.(Parsable)
		if ok {
			res = append(res, par)
			continue
		}
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("[parser:newseq] arg at idx %d not of valid type: %T", i+1, v)
		}
		seq, ok := p.seqs[str]
		if !ok {
			return nil, fmt.Errorf("[parser:newseq] no seq found with name %v", str)
		}
		res = append(res, seq)
	}
	_, ok := p.seqs[n]
	if ok {
		fmt.Printf("[parser:newseq] overwriting sequence with name \"%v\"", n)
	}
	new := NewSequence(n, res...)
	return new, nil
}

func (p *Parser) Using(n string) Parsable {
	par, ok := p.seqs[n]
	if !ok {
		return &noopParser{}
	}
	return par
}

type noopParser struct{}

func (n *noopParser) Parse(_ string) (*Result, error) {
	return retErr("noop", fmt.Errorf("[noop] no op parser used"))
}

func (n noopParser) Name() string { return "noop" }
