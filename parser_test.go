package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tcs := []struct {
		desc     string
		input    string
		expected *pl.Pl
	}{
		{
			desc:  "function with reference arguments",
			input: `(a $.b[1]["c-1"] $[2]["d-d"][3] $.e[4][f])`,
			expected: must(pl.NewPl(
				must(pl.NewFn("a",
					must(pl.NewRef("b", 1, "c-1")),
					must(pl.NewRef(2, "d-d", 3)),
					must(pl.NewRef("e", 4, "f")),
				)),
			)),
		},
		{
			desc:  "function with multiple arguments",
			input: `(a "b" 42 $.a[1].b 3.14 "36")`,
			expected: must(pl.NewPl(
				must(pl.NewFn("a", "b", 42, must(pl.NewRef("a", 1, "b")), 3.14, "36")),
			)),
		},
		{
			desc:  "sequence of functions",
			input: `(a "b" 42 3.14 "36" | c "d" 21)`,
			expected: must(pl.NewPl(
				must(pl.NewFn("a", "b", 42, 3.14, "36")),
				must(pl.NewFn("c", "d", 21)),
			)),
		},
		{
			desc:  "nested function",
			input: `(a "b" (c "d" 21) 3.14 "36")`,
			expected: must(pl.NewPl(
				must(pl.NewFn("a", "b", must(pl.NewPl(must(pl.NewFn("c", "d", 21)))), 3.14, "36")))),
		},
		{
			desc:  "consecutive nested functions",
			input: `(a "b" (c "d" 21) (e 3.14) 37)`,
			expected: must(pl.NewPl(
				must(pl.NewFn(
					"a",
					"b",
					must(pl.NewPl(must(pl.NewFn("c", "d", 21)))),
					must(pl.NewPl(must(pl.NewFn("e", 3.14)))),
					37,
				)),
			)),
		},
		{
			desc:  "nested sequence of functions",
			input: `(a "b" (c "d" 21 | e 3.14) 37)`,
			expected: must(pl.NewPl(
				must(pl.NewFn(
					"a",
					"b",
					must(pl.NewPl(
						must(pl.NewFn("c", "d", 21)),
						must(pl.NewFn("e", 3.14)),
					)),
					37,
				)),
			)),
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			fns, err := pl.ParseString(tc.input)
			require.NoError(err)
			require.Equal(tc.expected, fns)
		})
	}
}
