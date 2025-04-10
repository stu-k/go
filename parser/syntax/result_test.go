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
	x := stx.NewResult("x", nil, "")
	y := "y"
	wantY := ss(y, y)
	yr := stx.NewResult(y, wantY, "")
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

func TestResultFor(t *testing.T) {
	res := stx.NewResult("abc", nil, "")
	res1 := stx.NewResult("123", nil, "")
	res.Append(res1)

	t.Logf("[abc] map: %v", res.ResultMap())
	if !res.HasResult("123") {
		t.Errorf("[abc] expected key 123 not to be empty")
	}

	r1 := res.ResultFor("xyz")
	if !r1.IsEmpy() {
		t.Errorf("[abc] expected key xyz to be empty; got %v", r1)
	}
}
