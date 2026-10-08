package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func isEven(i int) bool { return i%2 == 0 }

func TestMap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"nil input gives a non-nil empty slice", nil, []int{}},
		{"empty input gives a non-nil empty slice", []int{}, []int{}},
		{"maps each element in order", []int{1, 2, 3}, []int{2, 4, 6}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.in, func(i int) int { return i * 2 })
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFilter(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want []int
	}{
		{"nil input gives a non-nil empty slice", nil, []int{}},
		{"no match gives a non-nil empty slice", []int{1, 3}, []int{}},
		{"keeps matches in order", []int{1, 2, 3, 4}, []int{2, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.in, isEven)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want bool
	}{
		{"nil input is false", nil, false},
		{"no match is false", []int{1, 3}, false},
		{"one match is true", []int{1, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Any(tt.in, isEven))
		})
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{"nil input gives the zero value", nil, 0},
		{"no match gives the zero value", []int{1, 3}, 0},
		{"returns the first match", []int{1, 4, 2}, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Find(tt.in, isEven))
		})
	}
}

func TestFindOk(t *testing.T) {
	tests := []struct {
		name   string
		in     []int
		want   int
		wantOk bool
	}{
		{"nil input is not found", nil, 0, false},
		{"no match is not found", []int{1, 3}, 0, false},
		{"a non-zero match is found", []int{1, 4}, 4, true},
		// FindOk reports ok by comparing the result with the zero value, so a
		// match that is itself the zero value reads as not found. This pins the
		// current behaviour; callers that can match a zero value must not use it.
		{"a zero-valued match reads as not found", []int{1, 0}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := FindOk(tt.in, isEven)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantOk, ok)
		})
	}
}

func TestOrDefault(t *testing.T) {
	assert.Equal(t, "fallback", OrDefault("", "fallback"))
	assert.Equal(t, "set", OrDefault("set", "fallback"))
	assert.Equal(t, 7, OrDefault(0, 7))
	assert.Equal(t, -1, OrDefault(-1, 7))

	var nilPtr *int
	fallback := 3
	assert.Equal(t, &fallback, OrDefault(nilPtr, &fallback))
}

func TestCopySlice(t *testing.T) {
	assert.Equal(t, []int{}, CopySlice[int](nil))

	original := []int{1, 2}
	copied := CopySlice(original)
	copied[0] = 9
	assert.Equal(t, []int{1, 2}, original)
}
