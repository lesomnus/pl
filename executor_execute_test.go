package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestExecutorExecute(t *testing.T) {
	executor := pl.Executor{
		Funcs: map[string]any{
			"sum": func(vs ...int) int {
				rst := 0
				for _, v := range vs {
					rst += v
				}

				return rst
			},
			"twice": func(vs ...int) []int {
				for i, v := range vs {
					vs[i] = v * 2
				}

				return vs
			},
		},
	}

	tcs := []struct {
		desc     string
		pl       *pl.Pl
		expected []any
	}{
		{
			desc: "single",
			pl: &pl.Pl{Funcs: []*pl.Fn{
				must(pl.NewFn("sum", 1, 2, 3, 4, 5)),
			}},
			expected: []any{1 + 2 + 3 + 4 + 5},
		},
		{
			desc: "piped",
			pl: &pl.Pl{Funcs: []*pl.Fn{
				must(pl.NewFn("sum", 1, 2, 3, 4, 5)),
				must(pl.NewFn("sum", 6, 7)),
			}},
			expected: []any{1 + 2 + 3 + 4 + 5 + 6 + 7},
		},
		{
			desc: "nested",
			pl: &pl.Pl{Funcs: []*pl.Fn{
				must(pl.NewFn("sum",
					1, 2,
					&pl.Pl{Funcs: []*pl.Fn{must(pl.NewFn("sum", 6, 7))}},
					4, 5,
				)),
			}},
			expected: []any{1 + 2 + (6 + 7) + 4 + 5},
		},
		{
			desc: "nested function that returns multiple values",
			pl: &pl.Pl{Funcs: []*pl.Fn{
				must(pl.NewFn("sum",
					1, 2,
					&pl.Pl{Funcs: []*pl.Fn{must(pl.NewFn("twice", 6, 7))}},
					4, 5,
				)),
			}},
			expected: []any{1 + 2 + (6 * 2) + (7 * 2) + 4 + 5},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			rst, err := executor.Execute(tc.pl, nil)
			require.NoError(err)
			require.ElementsMatch(tc.expected, rst)
		})
	}

	t.Run("fails if", func(t *testing.T) {
		tcs := []struct {
			desc string
			pl   *pl.Pl
			msgs []string
		}{
			{
				desc: "function is not defined",
				pl:   &pl.Pl{Funcs: []*pl.Fn{{Name: "quantum carburetor", Args: []*pl.Arg{}}}},
				msgs: []string{"not defined"},
			},
			{
				desc: "argument is invalid",
				pl:   &pl.Pl{Funcs: []*pl.Fn{{Name: "sum", Args: []*pl.Arg{{Int: addr(42)}, {}}}}},
				msgs: []string{"fn[0]", "sum", "arg[1]"},
			},
			{
				desc: "nested function is failed when evaluate",
				pl:   &pl.Pl{Funcs: []*pl.Fn{{Name: "sum", Args: must(pl.NewArgs(42, &pl.Pl{Funcs: []*pl.Fn{{Name: "Jerry Smith (C-131)", Args: []*pl.Arg{}}}}))}}},
				msgs: []string{"fn[0]", "sum", "arg[1]"},
			},
			{
				desc: "nested function is failed when invoke",
				pl:   &pl.Pl{Funcs: []*pl.Fn{{Name: "sum", Args: must(pl.NewArgs(42, &pl.Pl{Funcs: []*pl.Fn{must(pl.NewFn("sum", "Unity"))}}))}}},
				msgs: []string{"fn[0]", "sum", "arg[1]"},
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {
				require := require.New(t)

				_, err := executor.Execute(tc.pl, nil)
				for _, msg := range tc.msgs {
					require.ErrorContains(err, msg)
				}
			})
		}
	})
}

func TestExecutorExecuteExpr(t *testing.T) {
	executor := pl.Executor{
		Funcs: map[string]any{
			"sum": func(vs ...int) int {
				rst := 0
				for _, v := range vs {
					rst += v
				}

				return rst
			},
		},
	}

	rst, err := executor.ExecuteExpr(`(sum 1 2 (sum 3 | sum (sum 4 5) 6) 7 (sum 8) | sum 9 10)`, nil)
	require.NoError(t, err)
	require.Equal(t, []any{1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10}, rst)

	t.Run("fails if expression is invalid", func(t *testing.T) {
		require := require.New(t)

		executor := pl.NewExecutor()
		_, err := executor.ExecuteExpr("((sum 1 2))", nil)
		require.ErrorContains(err, "unexpected token")
	})
}
