package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func addr[T any](v T) *T {
	return &v
}

func must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

func TestNewRef(t *testing.T) {
	tcs := []struct {
		desc     string
		input    []interface{}
		expected []pl.RefKey
	}{
		{
			desc:     "keys",
			input:    []interface{}{"a", "b", "c"},
			expected: []pl.RefKey{{Name: addr("a")}, {Name: addr("b")}, {Name: addr("c")}},
		},
		{
			desc:  "indexes",
			input: []interface{}{1, -2, 3},
			expected: []pl.RefKey{
				{Index: addr(1)},
				{Index: addr(-2)},
				{Index: addr(3)},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			args, err := pl.NewRef(tc.input...)
			require.NoError(err)
			require.Equal(tc.expected, args)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require := require.New(t)

		_, err := pl.NewRef("baz", 3.14)
		require.ErrorContains(err, "invalid")
	})
}

func TestNewArgs(t *testing.T) {
	tcs := []struct {
		desc     string
		input    []interface{}
		expected []*pl.Arg
	}{
		{
			desc:  "scalars",
			input: []interface{}{"Rick", 42, 3.14, "36"},
			expected: []*pl.Arg{
				{String: addr("Rick")},
				{Int: addr(42)},
				{Float: addr(3.14)},
				{String: addr("36")},
			},
		},
		{
			desc:     "reference",
			input:    []interface{}{[]pl.RefKey{{Name: addr("b")}, {Index: addr(1)}, {Name: addr("c")}}},
			expected: []*pl.Arg{{Ref: []pl.RefKey{{Name: addr("b")}, {Index: addr(1)}, {Name: addr("c")}}}},
		},
		{
			desc: "nested function",
			input: []interface{}{&pl.Pl{
				Funcs: []*pl.Fn{
					{
						Name: "Zeep",
						Args: []*pl.Arg{
							{String: addr("Rick")},
							{Int: addr(42)},
							{Float: addr(3.14)},
							{String: addr("36")},
						},
					},
				},
			}},
			expected: []*pl.Arg{
				{Nested: &pl.Pl{
					Funcs: []*pl.Fn{
						{
							Name: "Zeep",
							Args: []*pl.Arg{
								{String: addr("Rick")},
								{Int: addr(42)},
								{Float: addr(3.14)},
								{String: addr("36")},
							},
						},
					},
				}},
			},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			args, err := pl.NewArgs(tc.input...)
			require.NoError(err)
			require.Equal(tc.expected, args)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require := require.New(t)

		_, err := pl.NewArgs(t)
		require.ErrorContains(err, "invalid")
	})
}

func TestNewFunc(t *testing.T) {
	tcs := []struct {
		desc     string
		name     string
		args     []interface{}
		expected *pl.Fn
	}{
		{
			desc:     "without args",
			name:     "Dio",
			args:     []interface{}{},
			expected: &pl.Fn{Name: "Dio", Args: []*pl.Arg{}},
		},
		{
			desc: "with args",
			name: "JoJo",
			args: []interface{}{"Jotaro", "Rohan"},
			expected: &pl.Fn{Name: "JoJo", Args: []*pl.Arg{
				{String: addr("Jotaro")},
				{String: addr("Rohan")},
			}},
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			f, err := pl.NewFn(tc.name, tc.args...)
			require.NoError(err)
			require.Equal(tc.expected, f)
		})
	}

	t.Run("invalid type", func(t *testing.T) {
		require := require.New(t)

		_, err := pl.NewFn("Joseph", t)
		require.ErrorContains(err, "invalid")
	})
}

func TestNewPl(t *testing.T) {
	require := require.New(t)

	f, err := pl.NewFn("a")
	require.NoError(err)

	plan, err := pl.NewPl(f)
	require.NoError(err)
	require.Equal(&pl.Pl{Funcs: []*pl.Fn{f}}, plan)
}
