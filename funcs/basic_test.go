package funcs_test

import (
	"testing"

	"github.com/lesomnus/pl/funcs"
	"github.com/stretchr/testify/require"
)

func TestPass(t *testing.T) {
	require := require.New(t)

	input := []any{"pi", 3.14, "answer", 42}
	output := funcs.Pass(input...)

	require.Equal(input, output)
}

func TestPrintf(t *testing.T) {
	require := require.New(t)

	output := funcs.Printf("%s %d %.2f", "a", 42, 3.14)
	require.Equal(output, "a 42 3.14")
}
