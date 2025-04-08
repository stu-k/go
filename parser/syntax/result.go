package syntax

type ParseResult interface {
	Name() string
	Strings() []string
	Rest() string
}

type parseResult struct {
	name    string
	results []string
	rest    string
}

func NewParseResult(n string, s []string, r string) *parseResult {
	if s == nil {
		s = []string{}
	}
	return &parseResult{n, s, r}
}

func (p *parseResult) Name() string {
	return p.name
}

func (p *parseResult) Strings() []string {
	return p.results
}

func (p *parseResult) Rest() string {
	return p.rest
}

func (p *parseResult) Append(r ParseResult) {
	if r == nil {
		return
	}
	p.results = append(p.results, r.Strings()...)
}

func (p *parseResult) SetRest(r string) {
	p.rest = r
}

func (p *parseResult) Len() int {
	if p.results == nil {
		return 0
	}
	return len(p.Strings())
}
