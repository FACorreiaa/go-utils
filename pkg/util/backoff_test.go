package util

import (
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewBackoff(t *testing.T) {
	t.Run("BackoffStrategyStatic", func(t *testing.T) {
		backoff := NewBackoff(BackoffStrategyStatic(time.Minute))
		assert.Equal(t, time.Minute, backoff(1))
	})
	t.Run("BackoffStrategyExponential", func(t *testing.T) {
		backoff := NewBackoff(BackoffStrategyExponential(time.Minute, 2))
		assert.Equal(t, time.Minute*1, backoff(0))
		assert.Equal(t, time.Minute*4, backoff(2))
		assert.Equal(t, time.Minute*32, backoff(5))
	})
	t.Run("BackoffWithMax should limit to max", func(t *testing.T) {
		backoff := NewBackoff(BackoffWithMax(time.Second*30), BackoffStrategyStatic(time.Second*10))
		assert.Equal(t, time.Second*10, backoff(1))
	})
	t.Run("BackoffWithMax should not limit to max", func(t *testing.T) {
		backoff := NewBackoff(BackoffWithMax(time.Second*30), BackoffStrategyStatic(time.Minute))
		assert.Equal(t, time.Second*30, backoff(1))
	})
	t.Run("no strategy is a zero delay, not a panic", func(t *testing.T) {
		backoff := NewBackoff(BackoffWithMax(time.Second))
		assert.Equal(t, time.Duration(0), backoff(3))
	})
	t.Run("exponential saturates instead of overflowing", func(t *testing.T) {
		backoff := NewBackoff(BackoffStrategyExponential(time.Second, 2))
		for _, attempt := range []int{33, 34, 40, 64, 1000} {
			assert.Positive(t, backoff(attempt), "attempt %d", attempt)
		}
		assert.Equal(t, time.Duration(math.MaxInt64), backoff(1000))
	})
	t.Run("exponential with max stays at max for large attempts", func(t *testing.T) {
		backoff := NewBackoff(BackoffStrategyExponential(time.Second, 2), BackoffWithMax(time.Minute))
		for _, attempt := range []int{6, 34, 100} {
			assert.Equal(t, time.Minute, backoff(attempt), "attempt %d", attempt)
		}
	})
	t.Run("jitter cannot push a saturated delay past max int", func(t *testing.T) {
		backoff := NewBackoff(BackoffStrategyExponential(time.Second, 2), BackoffWithJitter(time.Second))
		assert.Positive(t, backoff(1000))
	})
}
