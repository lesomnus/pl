package pl

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecutorEvaluateFn(t *testing.T) {
	data := map[string]any{
		"a": []map[string]string{{
			"b": "foo",
		}},
	}

	executor := Executor{}

	t.Run("resolve arguments from args", func(t *testing.T) {
		require := require.New(t)

		ref, err := NewRef("a", 0, "b")
		require.NoError(err)

		args, err := NewArgs("string", 3.14, 42, ref, &Pl{})
		require.NoError(err)

		fn := &Fn{Name: "fn", Args: args}

		node, err := executor.evaluateFn(fn, data)
		require.NoError(err)
		require.Equal("fn", node.name)
		require.ElementsMatch([]any{"string", 3.14, 42, "foo", &Pl{}}, node.args)
	})

	t.Run("fails if reference is not resolved", func(t *testing.T) {
		require := require.New(t)

		ref, err := NewRef("c", 1, "z")
		require.NoError(err)

		args, err := NewArgs("string", ref)
		require.NoError(err)

		fn := &Fn{Name: "fn", Args: args}

		_, err = executor.evaluateFn(fn, data)
		require.Error(err)
		require.ErrorContains(err, "arg[1]")
		require.ErrorContains(err, "reference")
	})

	t.Run("fails if argument is empty", func(t *testing.T) {
		require := require.New(t)

		args, err := NewArgs("string", 42)
		require.NoError(err)

		args = append(args, &Arg{})

		fn := &Fn{Name: "fn", Args: args}

		_, err = executor.evaluateFn(fn, data)
		require.Error(err)
		require.ErrorContains(err, "arg[2]")
		require.ErrorContains(err, "empty")
	})
}

func TestExecutorInvoke(t *testing.T) {
	executor := Executor{Convs: defaultConversions}
	executor.Convs.Set(string_t, reflect.TypeOf(float64(0)), func(v reflect.Value) (any, error) {
		return strconv.ParseFloat(v.String(), 64)
	})
	executor.Convs.Set(reflect.TypeOf(float32(0)), string_t, func(v reflect.Value) (any, error) {
		return nil, errors.New("fail")
	})

	sum := func(vs ...int) int {
		rst := 0
		for _, v := range vs {
			rst += v
		}

		return rst
	}

	cat := func(vs ...string) string {
		rst := ""
		for _, v := range vs {
			rst += v
		}

		return rst
	}

	tcs := []struct {
		desc string
		fn   any
		args []any
		rst  any
	}{
		{
			desc: "invoke a function without arguments",
			fn:   func() int { return 42 },
			args: []any{},
			rst:  42,
		},
		{
			desc: "invoke a function with argument",
			fn:   func(v int) int { return v * 2 },
			args: []any{17},
			rst:  34,
		},
		{
			desc: "invoke a function with multiple arguments",
			fn:   func(lhs int, rhs int) int { return lhs + rhs },
			args: []any{19, 36},
			rst:  55,
		},
		{
			desc: "invoke a variadic function without arguments",
			fn:   sum,
			args: []any{},
			rst:  0,
		},
		{
			desc: "invoke a variadic function with argument",
			fn:   sum,
			args: []any{42},
			rst:  42,
		},
		{
			desc: "invoke a variadic function with multiple arguments",
			fn:   sum,
			args: []any{1, 2, 3, 4, 5},
			rst:  15,
		},
		{
			desc: "invoke a variadic function with implicit conversion",
			fn:   cat,
			args: []any{42, " ", 31},
			rst:  "42 31",
		},
		{
			desc: "invoke a function with implicit conversion",
			fn:   func(v int) string { return strconv.Itoa(v) },
			args: []any{42},
			rst:  "42",
		},
		{
			desc: "invoke a function with implicit conversion to string if method (struct) String() by struct exists",
			fn:   func(v string) string { return v },
			args: []any{time.Date(1955, 11, 12, 22, 4, 0, 0, time.FixedZone("UTC-7", -7*50*50))},
			rst:  "1955-11-12 22:04:00 -0451 UTC-7",
		},
		{
			desc: "invoke a function with implicit conversion to string if method (struct) String() by pointer to struct exists",
			fn:   func(v string) string { return v },
			args: []any{(func() *time.Time {
				d := time.Date(1955, 11, 12, 22, 4, 0, 0, time.FixedZone("UTC-7", -7*50*50))
				return &d
			})()},
			rst: "1955-11-12 22:04:00 -0451 UTC-7",
		},
		{
			desc: "invoke a function with implicit conversion to string if method (*struct) String() by pointer to struct exists",
			fn:   func(v string) string { return v },
			args: []any{(func() *strings.Builder {
				sb := strings.Builder{}
				sb.WriteString("Josuke")
				sb.WriteString(" Higashikata")
				return &sb
			})()},
			rst: "Josuke Higashikata",
		},
		{
			desc: "invoke a function with intermediate conversion",
			fn:   func(v float64) string { return fmt.Sprintf("%.2f", v) },
			args: []any{(func() *strings.Builder {
				sb := strings.Builder{}
				sb.WriteString("3.14")
				return &sb
			})()},
			rst: "3.14",
		},
		{
			desc: "function can returns an error",
			fn:   func() (string, error) { return "Zoidberg", nil },
			args: []any{},
			rst:  "Zoidberg",
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			rst, err := executor.invokeFn(tc.fn, tc.args)
			require.NoError(err)
			require.Equal(tc.rst, rst)
		})
	}

	t.Run("fails if", func(t *testing.T) {
		tcs := []struct {
			desc string
			fn   any
			args []any
			msgs []string
		}{
			{
				desc: "return nothing",
				fn:   func() {},
				args: []any{},
				msgs: []string{"one or two"},
			},
			{
				desc: "return more than two values",
				fn:   func() (int, string, error) { return 42, "morty", nil },
				args: []any{},
				msgs: []string{"one or two"},
			},
			{
				desc: "return two values without error type",
				fn:   func() (int, string) { return 21, "rick" },
				args: []any{},
				msgs: []string{"error"},
			},
			{
				desc: "number of arguments not fit",
				fn:   func(int) int { return 41 },
				args: []any{31, 53},
				msgs: []string{"args are given"},
			},
			{
				desc: "number of arguments not fit to a variadic function",
				fn:   func(int, string, ...int) int { return 31 },
				args: []any{55},
				msgs: []string{"at least"},
			},
			{
				desc: "invalid type of argument",
				fn:   func(int, int) int { return 4 },
				args: []any{1055, "bender"},
				msgs: []string{"arg[1]"},
			},
			{
				desc: "function can returns an error",
				fn:   func() (string, error) { return "Hubert", errors.New("Farnsworth") },
				args: []any{},
				msgs: []string{"Farnsworth"},
			},
			{
				desc: "invoke a function with implicit conversion fails",
				fn:   func(v float64) string { return "" },
				args: []any{"pi"},
				msgs: []string{"arg[0]"},
			},
			{
				desc: "invoke a function with implicit conversion to string if method String() not exists",
				fn:   func(v string) string { return v },
				args: []any{errors.New("Cronenbergs")},
				msgs: []string{"arg[0]"},
			},
			{
				desc: "invoke a function with implicit conversion to string if method (*struct) String() by struct exists",
				fn:   func(v string) string { return v },
				args: []any{(func() strings.Builder {
					sb := strings.Builder{}
					sb.WriteString("Josuke")
					sb.WriteString(" Higashikata")
					return sb
				})()},
				msgs: []string{"not found"},
			},
			{
				// float32->string fails.
				desc: "invoke a function with intermediate conversion to fails",
				fn:   func(v float64) string { return "" },
				args: []any{float32(3.14)},
				msgs: []string{"to intermediate"},
			},
			{
				// strings.Builder->string OK, string->float64 fails.
				desc: "invoke a function with intermediate conversion from fails",
				fn:   func(v float64) string { return "" },
				args: []any{(func() *strings.Builder {
					sb := strings.Builder{}
					sb.WriteString("Pi")
					return &sb
				})()},
				msgs: []string{"from intermediate"},
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {
				require := require.New(t)

				_, err := executor.invokeFn(tc.fn, tc.args)
				require.Error(err)
				for _, msg := range tc.msgs {
					require.ErrorContains(err, msg)
				}
			})
		}
	})
}
