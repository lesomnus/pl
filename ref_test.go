package pl_test

import (
	"testing"

	"github.com/lesomnus/pl"
	"github.com/stretchr/testify/require"
)

func TestReString(t *testing.T) {
	require := require.New(t)

	ref := pl.Ref{
		{Name: addr("foo")},
		{Index: addr(42)},
		{},
	}

	path := ref.String()
	require.Equal(".foo[42].?", path)
}

func TestResolve(t *testing.T) {
	tcs := []struct {
		desc     string
		input    any
		ref      pl.Ref
		expected any
	}{
		{
			desc:     "map of string key",
			input:    map[string]int{"answer": 42},
			ref:      must(pl.NewRef("answer")),
			expected: 42,
		},
		{
			desc:     "pointer to map string key",
			input:    &map[string]int{"answer": 42},
			ref:      must(pl.NewRef("answer")),
			expected: 42,
		},
		{
			desc:     "map of int key",
			input:    map[int]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of int16 key",
			input:    map[int16]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of int32 key",
			input:    map[int32]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of int64 key",
			input:    map[int64]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of uint key",
			input:    map[uint]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of uint16 key",
			input:    map[uint16]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of uint32 key",
			input:    map[uint32]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "map of uint64 key",
			input:    map[uint64]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "pointer to map of int key",
			input:    map[int]string{42: "answer"},
			ref:      must(pl.NewRef(42)),
			expected: "answer",
		},
		{
			desc:     "struct",
			input:    struct{ Answer int }{Answer: 42},
			ref:      must(pl.NewRef("Answer")),
			expected: 42,
		},
		{
			desc:     "pointer to struct",
			input:    &struct{ Answer int }{Answer: 42},
			ref:      must(pl.NewRef("Answer")),
			expected: 42,
		},
		{
			desc:     "array",
			input:    [3]string{"foo", "bar", "baz"},
			ref:      must(pl.NewRef(1)),
			expected: "bar",
		},
		{
			desc:     "pointer to array",
			input:    &[3]string{"foo", "bar", "baz"},
			ref:      must(pl.NewRef(1)),
			expected: "bar",
		},
		{
			desc:     "slice",
			input:    []string{"foo", "bar", "baz"},
			ref:      must(pl.NewRef(1)),
			expected: "bar",
		},
		{
			desc:     "pointer to slice",
			input:    &[]string{"foo", "bar", "baz"},
			ref:      must(pl.NewRef(1)),
			expected: "bar",
		},
		{
			desc:     "nested",
			input:    struct{ A any }{A: []any{map[string]any{"b": map[int]any{42: &struct{ Pi float64 }{Pi: 3.14}}}}},
			ref:      must(pl.NewRef("A", 0, "b", 42, "Pi")),
			expected: 3.14,
		},
	}
	for _, tc := range tcs {
		t.Run(tc.desc, func(t *testing.T) {
			require := require.New(t)

			rst, err := pl.Resolve(tc.input, tc.ref)
			require.NoError(err)
			require.Equal(tc.expected, rst)
		})
	}

	t.Run("fails if", func(t *testing.T) {
		tcs := []struct {
			desc  string
			input any
			ref   pl.Ref
			msgs  []string
		}{
			{
				desc:  "map of string key with non-exist key",
				input: map[string]int{"answer": 42},
				ref:   must(pl.NewRef("pi")),
				msgs:  []string{"no", "key", "pi"},
			},
			{
				desc:  "map of int key with non-exist key",
				input: map[int]string{42: "answer"},
				ref:   must(pl.NewRef(36)),
				msgs:  []string{"no", "key", "36"},
			},
			{
				desc:  "map of non-string key with string",
				input: map[int]string{42: "answer"},
				ref:   must(pl.NewRef("answer")),
				msgs:  []string{"key", "not", "string"},
			},
			{
				desc:  "map of non-int key with int",
				input: map[string]int{"answer": 42},
				ref:   must(pl.NewRef(42)),
				msgs:  []string{"key", "not", "integer"},
			},
			{
				desc:  "struct with non-exist field",
				input: struct{ Answer int }{Answer: 42},
				ref:   must(pl.NewRef("Pi")),
				msgs:  []string{"no", "field", "Pi"},
			},
			{
				desc:  "non-object using key",
				input: []string{},
				ref:   must(pl.NewRef("answer")),
				msgs:  []string{"not", "object"},
			},
			{
				desc:  "non-list using index",
				input: struct{}{},
				ref:   must(pl.NewRef(0)),
				msgs:  []string{"not", "list"},
			},
			{
				desc:  "out of range",
				input: []string{"foo", "bar", "baz"},
				ref:   must(pl.NewRef(3)),
				msgs:  []string{"out of range"},
			},
			{
				desc:  "invalid key",
				input: struct{}{},
				ref:   []pl.RefKey{{}},
				msgs:  []string{"invalid", "key"},
			},
		}
		for _, tc := range tcs {
			t.Run(tc.desc, func(t *testing.T) {
				require := require.New(t)

				_, err := pl.Resolve(tc.input, tc.ref)
				for _, msg := range tc.msgs {
					require.ErrorContains(err, msg)
				}
			})
		}
	})
}
