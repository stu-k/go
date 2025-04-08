package syntax_test

import (
	"testing"

	stx "github.com/stu-k/go/parser/syntax"
)

var meq = func(a, b map[string][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k1, v1 := range a {
		v2, ok := b[k1]
		if !ok {
			return false
		}
		if !eq(v1, v2) {
			return false
		}
	}
	return true
}

func TestResult(t *testing.T) {
	x := stx.NewParseResult("x", nil, "")
	y := "y"
	wantY := ss(y, y)
	yr := stx.NewParseResult(y, wantY, "")
	x.Append(yr)
	m := x.NameMap()

	gotY, ok := m[y]
	if !ok {
		t.Fatalf("expected key \"%v\" on result okMap", y)
	}
	if !eq(wantY, gotY) {
		t.Fatalf("expected okMap for key \"%v\" to be %v; got %v", y, wantY, gotY)
	}
}
