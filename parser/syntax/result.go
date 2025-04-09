package syntax

type ParseResult struct {
	name        string
	resultsStrs []string
	rest        string
	resultMap   map[string]*ParseResult
}

func NewParseResult(n string, s []string, r string) *ParseResult {
	if s == nil {
		s = []string{}
	}
	return &ParseResult{
		name:        n,
		resultsStrs: s,
		rest:        r,
		resultMap:   make(map[string]*ParseResult),
	}
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
	p.resultMap[r.Name()] = r
	p.resultsStrs = append(p.resultsStrs, r.Strings()...)
}

func (p *ParseResult) SetRest(r string) {
	p.rest = r
}

func (p *ParseResult) Len() int {
	if p.resultsStrs == nil || p.resultMap == nil {
		return 0
	}
	return len(p.resultMap)
}

func (p *ParseResult) NameMap() map[string][]string {
	m := make(map[string][]string)
	for _, result := range p.resultMap {
		m[result.Name()] = result.Strings()
	}
	return m
}

func (p *ParseResult) IsEmpy() bool {
	return len(p.resultsStrs) == 0 &&
		len(p.rest) == 0 &&
		len(p.resultMap) == 0
}

func (p *ParseResult) HasResult(name string) bool {
	_, ok := p.resultMap[name]
	return ok
}

func (p *ParseResult) ResultFor(name string) *ParseResult {
	r, ok := p.resultMap[name]
	if !ok {
		return NewParseResult("", nil, "")
	}
	return r
}

func (p *ParseResult) ResultMap() map[string]*ParseResult {
	m := make(map[string]*ParseResult)
	for _, result := range p.resultMap {
		m[result.Name()] = result
	}
	return m
}
