package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Elias-Chairi/ccnp-project-gruppe3/internal/util"
)

// -------------------------------------- Positive tests --------------------------------------
func TestStack_PushPopTop(t *testing.T) {
	assert := assert.New(t)
	var stack util.Stack[string]

	stack.Push("first")
	stack.Push("second")

	popped, err := stack.Pop()
	assert.NoError(err)
	assert.Equal("second", *popped)

	popped, err = stack.Pop()
	assert.NoError(err)
	assert.Equal("first", *popped)
}

// -------------------------------------- Negative tests --------------------------------------

func TestStack_PopEmpty(t *testing.T) {
	var stack util.Stack[float64]
	_, err := stack.Pop()
	assert.ErrorIs(t, err, util.ErrStackEmpty)
}
