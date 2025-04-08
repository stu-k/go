package syntax

type ParseResult struct {
	name        string
	resultsStrs []string
	rest        string
	results     []*ParseResult
}

func NewParseResult(n string, s []string, r string) *ParseResult {
	if s == nil {
		s = []string{}
	}
	return &ParseResult{n, s, r, nil}
}

func (p *ParseResult) Name() string {
	return p.name
}

func (p *ParseResult) Strings() []string {
	return p.resultsStrs
}

func (p *ParseResult) Rest() string {
	return p.rest
}

func (p *ParseResult) Append(r *ParseResult) {
	if r == nil {
		return
	}
	p.results = append(p.results, r)
	p.resultsStrs = append(p.resultsStrs, r.Strings()...)
}

func (p *ParseResult) SetRest(r string) {
	p.rest = r
}

func (p *ParseResult) Len() int {
	if p.resultsStrs == nil {
		return 0
	}
	return len(p.Strings())
}

func (p *ParseResult) NameMap() map[string][]string {
	m := make(map[string][]string)
	for _, result := range p.results {
		m[result.Name()] = result.Strings()
	}
	return m
}
