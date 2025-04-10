package syntax

type Result struct {
	name        string
	resultsStrs []string
	rest        string
	resultMap   map[string]*Result
}

func NewResult(name string, s []string, rest string) *Result {
	if s == nil {
		s = []string{}
	}
	return &Result{
		name:        name,
		resultsStrs: s,
		rest:        rest,
		resultMap:   make(map[string]*Result),
	}
}

func (p *Result) Name() string {
	return p.name
}

func (p *Result) Strings() []string {
	return p.resultsStrs
}

func (p *Result) Rest() string {
	return p.rest
}

func (p *Result) Append(r *Result) {
	if r == nil {
		return
	}
	p.resultMap[r.Name()] = r
	p.resultsStrs = append(p.resultsStrs, r.Strings()...)
}

func (p *Result) SetRest(r string) {
	p.rest = r
}

func (p *Result) Len() int {
	if p.resultsStrs == nil || p.resultMap == nil {
		return 0
	}
	return len(p.resultMap)
}

func (p *Result) NameMap() map[string][]string {
	m := make(map[string][]string)
	for _, result := range p.resultMap {
		m[result.Name()] = result.Strings()
	}
	return m
}

func (p *Result) IsEmpy() bool {
	return len(p.resultsStrs) == 0 &&
		len(p.rest) == 0 &&
		len(p.resultMap) == 0
}

func (p *Result) HasResult(name string) bool {
	_, ok := p.resultMap[name]
	return ok
}

func (p *Result) ResultFor(name string) *Result {
	r, ok := p.resultMap[name]
	if !ok {
		return NewResult("", nil, "")
	}
	return r
}

func (p *Result) ResultMap() map[string]*Result {
	m := make(map[string]*Result)
	for _, result := range p.resultMap {
		m[result.Name()] = result
	}
	return m
}
