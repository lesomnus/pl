package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
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
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			rst, err := executor.Execute(pl.NewPl(tc.input))
			require.NoError(err)
			require.Equal(tc.expected, rst)
		})
	}
}
