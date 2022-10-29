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
	Ref    []RefKey `parser:"| '$' @@+"`
	Nested *Pl      `parser:"| @@"`
}

type RefKey struct {
	Name  *string `parser:"  (('.' @(Ident|String)) | ('[' @(Ident|String) ']'))"`
	Index *int    `parser:"| '[' @Int ']'"`
}

func (k *RefKey) String() string {
	if k.Name != nil {
		return fmt.Sprintf(".%s", *k.Name)
	} else {
		return fmt.Sprintf("[%d]", *k.Index)
	}
}

func NewRefKey[T string | int](key T) RefKey {
	switch v := interface{}(key).(type) {
	case string:
		return RefKey{Name: &v}
	case int:
		return RefKey{Index: &v}
	}

	panic(key)
}

func K[T string | int](key T) RefKey {
	return NewRefKey(key)
}

var plParser = participle.MustBuild[Pl](
	participle.Unquote("String"),
)

func ParseString(expr string) (*Pl, error) {
	return plParser.ParseString("", expr)
}
