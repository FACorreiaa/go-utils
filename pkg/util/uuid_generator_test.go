package util

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _ UUIDGenerator = (*uuidGenerator)(nil)

func uuidVersion(u uuid.UUID) byte { return u[6] >> 4 }

// isRFC9562Variant reports whether the two most significant bits of octet 8 are 10.
func isRFC9562Variant(u uuid.UUID) bool { return u[8]&0xc0 == 0x80 }

func TestUUIDGenerator(t *testing.T) {
	assert.NotEmpty(t, NewUUIDGenerator().Generate())
}

func TestUUIDGenerator_Generate(t *testing.T) {
	g := NewUUIDGenerator()

	s := g.Generate()
	parsed, err := uuid.Parse(s)
	require.NoError(t, err)
	assert.Equal(t, s, parsed.String())
	assert.Len(t, s, 36)
	assert.EqualValues(t, 4, uuidVersion(parsed))
	assert.True(t, isRFC9562Variant(parsed))
	assert.NotEqual(t, s, g.Generate())
}

func TestUUIDGenerator_GenerateUUID(t *testing.T) {
	g := NewUUIDGenerator()

	u := g.GenerateUUID()
	assert.NotEqual(t, uuid.Nil(), u)
	assert.EqualValues(t, 4, uuidVersion(u))
	assert.True(t, isRFC9562Variant(u))

	seen := make(map[uuid.UUID]struct{}, 1000)
	for range 1000 {
		seen[g.GenerateUUID()] = struct{}{}
	}
	assert.Len(t, seen, 1000)
}

func TestUUIDv5(t *testing.T) {
	// Expected values from Python's uuid.uuid5(uuid.UUID(int=0), input), which
	// also match the previous github.com/google/uuid implementation.
	tests := []struct {
		input string
		want  string
	}{
		{"", "e129f27c-5103-5c5c-844b-cdf0a15e160d"},
		{"hello", "b7502f40-1152-59f2-ba10-69aeed522cdf"},
		{"Hello", "5e36032f-fc05-5506-8dff-e29438639449"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := UUIDv5(tt.input)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, got, UUIDv5(tt.input))

			parsed, err := uuid.Parse(got)
			require.NoError(t, err)
			assert.EqualValues(t, 5, uuidVersion(parsed))
			assert.True(t, isRFC9562Variant(parsed))
		})
	}
}

func TestNewSHA1(t *testing.T) {
	// RFC 9562 Appendix A.4 test vector.
	dns := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	assert.Equal(t, "2ed6657d-e927-568b-95e1-2665a8aea6a2", newSHA1(dns, []byte("www.example.com")).String())

	// Same name in a different namespace yields a different UUID.
	assert.NotEqual(t, newSHA1(uuid.Nil(), []byte("www.example.com")), newSHA1(dns, []byte("www.example.com")))
}
