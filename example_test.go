package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestExample(t *testing.T) {
	require := require.New(t)

	executor := pl.NewExecutor()
	executor.Funcs["sum"] = func(vs ...int) int {
		rst := 0
		for _, v := range vs {
			rst += v
		}

		return rst
	}

	rst, err := executor.ExecuteExpr("(sum 1 2 (sum 3 | sum (sum $.Answer 5) 6) 7 (sum 8) | sum 9 10)", struct{ Answer int }{Answer: 42})
	require.NoError(err)

	v, ok := rst[0].(int)
	require.True(ok)
	require.Equal(93, v)
}
