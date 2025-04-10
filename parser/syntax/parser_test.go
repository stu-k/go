package syntax_test

import (
	"fmt"
	"testing"
	"unicode"

	stx "github.com/stu-k/go/parser/syntax"
)

func TestParser(t *testing.T) {
	p := stx.NewParser("main")
	p.NewSeq("alpha", stx.NewRule("alpha").CheckChar(unicode.IsLetter))
	p.NewSeq("num", stx.NewRule("num").CheckChar(unicode.IsNumber))
	p.NewSeq("apos", stx.NewRule("apos").Repeat(1).Chars("'").Capture(false))
	p.NewSeq("comma", stx.NewRule("comma").Repeat(1).Chars(",").Capture(false))
	p.NewSeq("rbracket", stx.NewRule("lbracket").Repeat(1).Chars("[").Capture(false))
	p.NewSeq("lbarcket", stx.NewRule("rbracket").Repeat(1).Chars("]").Capture(false))
	p.NewSeq("alphanum", "num", "alpha")
	t.Run("one", func(t *testing.T) {
		res, err := p.Using("alpha").Parse("abc123")
		fmt.Printf("res: %+v\nerr: %v\n", res, err)
	})
	t.Run("two", func(t *testing.T) {
		res, err := p.Using("alphanum").Parse("abc123")
		fmt.Printf("res: %+v\nerr: %v\n", res, err)
	})

	// p.NewPickOne("alphanum",
	// )
	// str := p.NewSeq("quoted",
	// 	apos, alphanum, apos,
	// )

	// anyVal := stx.NewSequence("anyval",
	// 	alp, num, str,
	// ).PickOne()
	// valComma := stx.NewSequence("anyval comma",
	// 	anyVal, comma,
	// )
	// valCommaEnd := stx.NewSequence("anyval comma end optional",
	// 	valComma.UntilFail(),
	// 	anyVal,
	// )

	// arr := stx.NewSequence("arr",
	// 	lbracket,
	// 	valCommaEnd,
	// 	rbracket,
	// )
}
