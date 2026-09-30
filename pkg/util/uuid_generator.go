package util

import (
	"crypto/sha1"
	"uuid"
)

type IDGenerator interface {
	Generate() string
}

type UUIDGenerator interface {
	IDGenerator
	GenerateUUID() uuid.UUID
}

type uuidGenerator struct{}

func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}

func (c *uuidGenerator) Generate() string {
	return uuid.New().String()
}

func (c *uuidGenerator) GenerateUUID() uuid.UUID {
	return uuid.New()
}

// UUIDv5 returns a deterministic name-based UUID (RFC 9562 §5.5) in the Nil
// namespace. The stdlib uuid package has no v5 constructor, so it is built here.
func UUIDv5(input string) string {
	return newSHA1(uuid.Nil(), []byte(input)).String()
}

func newSHA1(space uuid.UUID, name []byte) uuid.UUID {
	h := sha1.New()
	h.Write(space[:])
	h.Write(name)
	var u uuid.UUID
	copy(u[:], h.Sum(nil))
	u[6] = (u[6] & 0x0f) | 0x50 // version 5
	u[8] = (u[8] & 0x3f) | 0x80 // RFC 9562 variant
	return u
}
