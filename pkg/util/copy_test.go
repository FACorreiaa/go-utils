package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopy(t *testing.T) {
	type item struct{ Name string }

	assert.Nil(t, Copy[item](nil, func(i *item) { i.Name = "changed" }))

	original := &item{Name: "original"}
	copied := Copy(original, func(i *item) { i.Name = "changed" })
	assert.Equal(t, "original", original.Name)
	assert.Equal(t, "changed", copied.Name)
	assert.NotSame(t, original, copied)
}
