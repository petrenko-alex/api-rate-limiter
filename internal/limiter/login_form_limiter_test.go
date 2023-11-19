package limiter_test

import "testing"

func TestLoginFormLimiter_SatisfyLimit(t *testing.T) {
	t.Run("ip in white list", func(t *testing.T) {
		// satisfy limit
		// NOT satisfy limit
	})

	t.Run("ip in black list", func(t *testing.T) {
		// not satisfy limit
		// satisfy limit
	})

	t.Run("no ip in lists", func(t *testing.T) {
		// satisfy
		// not satisfy
	})
}

func TestLoginFormLimiter_SatisfyLimit_Error(t *testing.T) {
	t.Run("incorrect identity", func(t *testing.T) {

	})
}
