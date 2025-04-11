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
	p.NewSeq("alphanum", stx.NewRule("alphanum").CheckChar(func(r rune) bool {
		return unicode.IsLetter(r) || unicode.IsNumber(r)
	}))
	p.NewSeq("quoted", "apos", "alphanum", "apos")
	p.NewPickOne("val", "alpha", "num", "quoted")
	p.NewUntilFail("val,", "val", "comma")
	p.NewSeq("val,val", "val,", "val")
	p.NewSeq("arr", "lbracket", "val,val", "rbracket")

	t.Run("arr", func(t *testing.T) {
		got, _ := p.Using("arr").Parse("[1,'2','a',b]")
		fmt.Printf("got: %+v", got)
	})
}
