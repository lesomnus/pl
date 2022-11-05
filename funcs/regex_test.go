package funcs_test

import (
	"testing"

	"github.com/lesomnus/pl/funcs"
	"github.com/stretchr/testify/require"
)

func TestRegex(t *testing.T) {
	t.Run("find match by index and name", func(t *testing.T) {
		require := require.New(t)

		matched, err := funcs.Regex(
			`([a-zA-Z]+) (?P<last>[a-zA-Z]+)`,
			"Seunghyun Hwang",
			"Patrick",
		)
		require.NoError(err)
		require.Len(matched, 1) // Patrick is not matched
		require.Equal("Seunghyun Hwang", matched[0].Source)
		require.Equal("Seunghyun Hwang", matched[0].String())
		require.Equal("Seunghyun", matched[0].ByIndex[0])
		require.Equal("Hwang", matched[0].ByIndex[1])
		require.Equal("Hwang", matched[0].ByName["last"])
	})

	t.Run("fails if invalid regex expression", func(t *testing.T) {
		require := require.New(t)

		_, err := funcs.Regex(`(`)
		require.ErrorContains(err, "compile")
	})
}
