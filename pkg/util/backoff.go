package util

import (
	"math"
	"math/rand"
	"time"
)

type Backoff func(attempt int) time.Duration

type backoffOptions struct {
	limit    time.Duration
	jitter   time.Duration
	strategy func(attempt int) time.Duration
}

type BackoffOption interface {
	apply(*backoffOptions)
}

type backoffMaxOption time.Duration

func (o backoffMaxOption) apply(opts *backoffOptions) {
	opts.limit = time.Duration(o)
}

func BackoffWithMax(max time.Duration) BackoffOption {
	return backoffMaxOption(max)
}

// Jitter backoff
type backoffJitterOption time.Duration

func (o backoffJitterOption) apply(opts *backoffOptions) {
	opts.jitter = time.Duration(o)
}

func BackoffWithJitter(jitter time.Duration) BackoffOption {
	return backoffJitterOption(jitter)
}

// Static backoff
type backoffStrategyStaticOption time.Duration

func (o backoffStrategyStaticOption) apply(opts *backoffOptions) {
	opts.strategy = func(_ int) time.Duration {
		return time.Duration(o)
	}
}

func BackoffStrategyStatic(static time.Duration) BackoffOption {
	return backoffStrategyStaticOption(static)
}

// Exponential backoff
type backoffStrategyExponentialOption struct {
	initial time.Duration
	base    float64
}

func (o backoffStrategyExponentialOption) apply(opts *backoffOptions) {
	opts.strategy = func(attempt int) time.Duration {
		factor := math.Pow(o.base, float64(attempt))
		// Saturate instead of overflowing: past this point the product wraps
		// negative and the caller would retry with no wait at all.
		if o.initial > 0 && factor >= float64(math.MaxInt64)/float64(o.initial) {
			return time.Duration(math.MaxInt64)
		}
		return o.initial * time.Duration(factor)
	}
}

func BackoffStrategyExponential(initial time.Duration, base int) BackoffOption {
	return backoffStrategyExponentialOption{initial: initial, base: float64(base)}
}

// NewBackoff builds a Backoff from opts. Without a strategy option the delay
// is zero (plus any jitter), rather than a nil-func panic on first use.
func NewBackoff(opts ...BackoffOption) Backoff {
	options := &backoffOptions{
		strategy: func(int) time.Duration { return 0 },
	}

	for _, o := range opts {
		o.apply(options)
	}

	return backoffWithOptions(options)
}

func backoffWithOptions(options *backoffOptions) Backoff {
	return func(attempt int) time.Duration {
		backoff := options.strategy(attempt)
		if options.jitter > 0 {
			jitter := time.Duration(rand.Int63n(int64(options.jitter)))
			if backoff > time.Duration(math.MaxInt64)-jitter {
				backoff = time.Duration(math.MaxInt64)
			} else {
				backoff += jitter
			}
		}
		if options.limit != 0 && backoff > options.limit {
			backoff = options.limit
		}
		return backoff
	}
}
