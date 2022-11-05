package pl

import (
	"fmt"

	"github.com/alecthomas/participle/v2"
)

type Pl struct {
	Funcs []*Fn `parser:"'(' ( @@ ( '|' @@ )* )? ')'"`
}

type Fn struct {
	Name string `parser:"@Ident"`
	Args []*Arg `parser:"@@*"`
}

type Arg struct {
	String *string  `parser:"  @String"`
	Float  *float64 `parser:"| @Float"`
	Int    *int     `parser:"| @Int"`
	Ref    Ref      `parser:"| '$' @@+"`
	Nested *Pl      `parser:"| @@"`
}

type RefKey struct {
	Name  *string `parser:"  (('.' @(Ident|String)) | ('[' @(Ident|String) ']'))"`
	Index *int    `parser:"| '[' @Int ']'"`
}

func (k *RefKey) String() string {
	if k.Name != nil {
		return fmt.Sprintf(".%s", *k.Name)
	} else if k.Index != nil {
		return fmt.Sprintf("[%d]", *k.Index)
	} else {
		return ".?"
	}
}

var plParser = participle.MustBuild[Pl](
	participle.Unquote("String"),
)

func ParseString(expr string) (*Pl, error) {
	return plParser.ParseString("", expr)
}
