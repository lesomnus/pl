package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/lesomnus/pl/funcs"
	"github.com/stretchr/testify/require"
)

func TestFuncMapDefaultFunctions(t *testing.T) {
	executor := pl.NewExecutor()

	tcs := []struct {
		desc     string
		input    *pl.Fn
		expected any
	}{
		{
			desc:     "pass",
			input:    &pl.Fn{Name: "pass", Args: must(pl.NewArgs("a", 42, 3.14))},
			expected: []any{"a", 42, 3.14},
		},
		{
			desc:     "printf",
			input:    &pl.Fn{Name: "printf", Args: must(pl.NewArgs("%s %d %.2f", "a", 42, 3.14))},
			expected: []any{"a 42 3.14"},
		},
		{
			desc:  "regex",
			input: &pl.Fn{Name: "regex", Args: must(pl.NewArgs(`foo([a-zA-Z]+)baz`, "foobarbaz"))},
			expected: []any{&funcs.RegexMatch{
				Source:  "foobarbaz",
				ByIndex: []string{"bar"},
				ByName:  make(map[string]string),
			}},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			rst, err := executor.Execute(pl.NewPl(tc.input), nil)
			require.NoError(err)
			require.Equal(tc.expected, rst)
		})
	}
}
