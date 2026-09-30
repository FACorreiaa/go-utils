package util

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
)

var (
	_ UUIDGenerator = (*MockUUIDGenerator)(nil)
	_ IDGenerator   = (*StaticUUIDGenerator)(nil)
)

func TestMockUUIDs(t *testing.T) {
	ids := []uuid.UUID{MockUUID1, MockUUID2, MockUUID3}
	for i, id := range ids {
		assert.NotEqual(t, uuid.Nil(), id)
		assert.EqualValues(t, 8, uuidVersion(id))
		assert.True(t, isRFC9562Variant(id))
		for _, other := range ids[i+1:] {
			assert.NotEqual(t, id, other)
		}
	}
	assert.Equal(t, "00000000-0000-8000-8000-000000000001", MockUUID1.String())
}

func TestMockUUIDGenerator(t *testing.T) {
	m := &MockUUIDGenerator{}
	m.On("Generate").Return("mock-id").Once()
	m.On("GenerateUUID").Return(MockUUID2).Once()

	assert.Equal(t, "mock-id", m.Generate())
	assert.Equal(t, MockUUID2, m.GenerateUUID())
	m.AssertExpectations(t)
}

func TestStaticUUIDGenerator(t *testing.T) {
	g := &StaticUUIDGenerator{}
	assert.Equal(t, "static-uuid", g.Generate())
	assert.Equal(t, g.Generate(), g.Generate())
}
