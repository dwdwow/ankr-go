package ankr

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Rate limiter implementation using token bucket algorithm.
//
// This package provides a rate limiter that can be used to limit
// the rate at which operations are performed.

// ============================================================================
// Rate Limit Behavior Types
// ============================================================================

// RateLimitBehavior defines how the rate limiter behaves when limit is exceeded
type RateLimitBehavior string

const (
	// RateLimitBlock waits until a token is available
	RateLimitBlock RateLimitBehavior = "block"
	// RateLimitRaise returns an error when rate limit is exceeded
	RateLimitRaise RateLimitBehavior = "raise"
	// RateLimitSkip returns false without executing when rate limit is exceeded
	RateLimitSkip RateLimitBehavior = "skip"
)

// ErrRateLimitExceeded is returned when rate limit is exceeded and behavior is RateLimitRaise
var ErrRateLimitExceeded = errors.New("rate limit exceeded")

// ============================================================================
// RateLimiter - Token Bucket Implementation
// ============================================================================

// RateLimiter implements the token bucket algorithm for rate limiting.
//
// This allows for burst traffic while maintaining an average rate limit over time.
type RateLimiter struct {
	limit           int64             // Maximum number of calls allowed in the period
	period          time.Duration     // Time period for the rate limit
	onLimitExceeded RateLimitBehavior // Behavior when rate limit is exceeded
	tokens          float64           // Current number of available tokens
	lastUpdate      time.Time         // Last time tokens were refilled
	mu              sync.RWMutex      // Mutex for thread-safe operations
}

// NewRateLimiter creates a new rate limiter with the specified parameters.
//
// Args:
//   - limit: Maximum number of calls allowed in the given period
//   - period: Time period for the rate limit
//   - onLimitExceeded: Behavior when rate limit is exceeded (default: RateLimitBlock)
//
// Returns:
//   - *RateLimiter: A new rate limiter instance
//   - error: Error if parameters are invalid
func NewRateLimiter(limit int64, period time.Duration, onLimitExceeded RateLimitBehavior) (*RateLimiter, error) {
	if limit <= 0 {
		return nil, errors.New("limit must be positive")
	}
	if period <= 0 {
		return nil, errors.New("period must be positive")
	}

	return &RateLimiter{
		limit:           limit,
		period:          period,
		onLimitExceeded: onLimitExceeded,
		tokens:          float64(limit),
		lastUpdate:      time.Now(),
	}, nil
}

// refillTokens refills tokens based on elapsed time since last refill.
func (rl *RateLimiter) refillTokens() {
	now := time.Now()
	elapsed := now.Sub(rl.lastUpdate)

	// Calculate tokens to add based on elapsed time
	tokensToAdd := elapsed.Seconds() * float64(rl.limit) / rl.period.Seconds()
	rl.tokens = min(float64(rl.limit), rl.tokens+tokensToAdd)
	rl.lastUpdate = now
}

// Acquire attempts to acquire tokens from the bucket.
//
// Args:
//   - ctx: Context for cancellation
//   - tokens: Number of tokens to acquire (default: 1)
//   - onLimitExceeded: Optional override for the default behavior
//
// Returns:
//   - bool: True if tokens were acquired, false otherwise
//   - error: Error if rate limit is exceeded and behavior is RateLimitRaise
func (rl *RateLimiter) Acquire(ctx context.Context, tokens int64, onLimitExceeded *RateLimitBehavior) (bool, error) {
	behavior := rl.onLimitExceeded
	if onLimitExceeded != nil {
		behavior = *onLimitExceeded
	}

	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.refillTokens()

	if rl.tokens >= float64(tokens) {
		rl.tokens -= float64(tokens)
		return true, nil
	}

	switch behavior {
	case RateLimitBlock:
		// Calculate wait time for next token
		tokensNeeded := float64(tokens) - rl.tokens
		waitTime := time.Duration(tokensNeeded * rl.period.Seconds() / float64(rl.limit) * float64(time.Second))

		// Release lock while waiting
		rl.mu.Unlock()

		// Wait with context support
		timer := time.NewTimer(waitTime)
		select {
		case <-timer.C:
			// Continue after wait
		case <-ctx.Done():
			timer.Stop()
			rl.mu.Lock()
			return false, ctx.Err()
		}

		// Re-acquire lock
		rl.mu.Lock()

		// Refill and try again
		rl.refillTokens()
		if rl.tokens >= float64(tokens) {
			rl.tokens -= float64(tokens)
			return true, nil
		}
		return false, nil

	case RateLimitRaise:
		return false, ErrRateLimitExceeded

	default: // RateLimitSkip
		return false, nil
	}
}

// AcquireWithTimeout attempts to acquire tokens with a timeout.
//
// This is a convenience method that creates a context with timeout.
func (rl *RateLimiter) AcquireWithTimeout(timeout time.Duration, tokens int64, onLimitExceeded *RateLimitBehavior) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return rl.Acquire(ctx, tokens, onLimitExceeded)
}

// TryAcquire attempts to acquire tokens without blocking.
//
// Returns true if tokens were acquired, false otherwise.
func (rl *RateLimiter) TryAcquire(tokens int64) bool {
	skip := RateLimitSkip
	acquired, _ := rl.Acquire(context.Background(), tokens, &skip)
	return acquired
}

// Reset resets the rate limiter to its initial state.
func (rl *RateLimiter) Reset() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	rl.tokens = float64(rl.limit)
	rl.lastUpdate = time.Now()
}

// GetAvailableTokens returns the current number of available tokens.
func (rl *RateLimiter) GetAvailableTokens() float64 {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	// Create a temporary copy to refill
	tmpRL := &RateLimiter{
		limit:      rl.limit,
		period:     rl.period,
		tokens:     rl.tokens,
		lastUpdate: rl.lastUpdate,
	}
	tmpRL.refillTokens()
	return tmpRL.tokens
}

// TimeUntilNextToken calculates time until next token is available.
func (rl *RateLimiter) TimeUntilNextToken() time.Duration {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	// Create a temporary copy to refill
	tmpRL := &RateLimiter{
		limit:      rl.limit,
		period:     rl.period,
		tokens:     rl.tokens,
		lastUpdate: rl.lastUpdate,
	}
	tmpRL.refillTokens()

	if tmpRL.tokens >= 1 {
		return 0
	}

	tokensNeeded := 1 - tmpRL.tokens
	waitSeconds := tokensNeeded * tmpRL.period.Seconds() / float64(tmpRL.limit)
	return time.Duration(waitSeconds * float64(time.Second))
}

// Wait blocks until a token is available or context is cancelled.
func (rl *RateLimiter) Wait(ctx context.Context) error {
	_, err := rl.Acquire(ctx, 1, nil)
	return err
}

// ============================================================================
// SharedRateLimiter - Shared Rate Limiter Across Multiple Operations
// ============================================================================

// SharedRateLimiter is a shared rate limiter that can be used across multiple operations.
//
// This is useful when you want to limit the combined rate of multiple
// operations rather than limiting each operation independently.
type SharedRateLimiter struct {
	limiter *RateLimiter
}

// NewSharedRateLimiter creates a new shared rate limiter.
func NewSharedRateLimiter(limit int64, period time.Duration, onLimitExceeded RateLimitBehavior) (*SharedRateLimiter, error) {
	limiter, err := NewRateLimiter(limit, period, onLimitExceeded)
	if err != nil {
		return nil, err
	}

	return &SharedRateLimiter{
		limiter: limiter,
	}, nil
}

// Acquire attempts to acquire tokens from the shared limiter.
func (srl *SharedRateLimiter) Acquire(ctx context.Context, tokens int64, onLimitExceeded *RateLimitBehavior) (bool, error) {
	return srl.limiter.Acquire(ctx, tokens, onLimitExceeded)
}

// TryAcquire attempts to acquire tokens without blocking.
func (srl *SharedRateLimiter) TryAcquire(tokens int64) bool {
	return srl.limiter.TryAcquire(tokens)
}

// Reset resets the rate limiter.
func (srl *SharedRateLimiter) Reset() {
	srl.limiter.Reset()
}

// GetAvailableTokens gets current available tokens.
func (srl *SharedRateLimiter) GetAvailableTokens() float64 {
	return srl.limiter.GetAvailableTokens()
}

// Wait blocks until a token is available.
func (srl *SharedRateLimiter) Wait(ctx context.Context) error {
	return srl.limiter.Wait(ctx)
}

// ============================================================================
// MultiRateLimiter - Multi-Tiered Rate Limiter
// ============================================================================

// RateLimit represents a single rate limit configuration.
type RateLimit struct {
	Limit  int64
	Period time.Duration
}

// MultiRateLimiter enforces multiple rate limits simultaneously.
//
// This is useful for APIs with tiered rate limiting, where you need to enforce
// multiple limits at different time scales (e.g., 10/second AND 100/minute AND 1000/hour).
// All rate limits must be satisfied for a request to proceed.
type MultiRateLimiter struct {
	limits          []RateLimit
	limiters        []*RateLimiter
	onLimitExceeded RateLimitBehavior
	mu              sync.RWMutex
}

// NewMultiRateLimiter creates a new multi-tiered rate limiter.
//
// Args:
//   - limits: List of rate limit configurations
//   - onLimitExceeded: Behavior when any rate limit is exceeded
//
// Example:
//
//	// 10 per second, 100 per minute, 1000 per hour
//	limiter, _ := NewMultiRateLimiter(
//	    []RateLimit{
//	        {Limit: 10, Period: time.Second},
//	        {Limit: 100, Period: time.Minute},
//	        {Limit: 1000, Period: time.Hour},
//	    },
//	    RateLimitBlock,
//	)
func NewMultiRateLimiter(limits []RateLimit, onLimitExceeded RateLimitBehavior) (*MultiRateLimiter, error) {
	if len(limits) == 0 {
		return nil, errors.New("at least one rate limit must be specified")
	}

	mrl := &MultiRateLimiter{
		limits:          limits,
		limiters:        make([]*RateLimiter, 0, len(limits)),
		onLimitExceeded: onLimitExceeded,
	}

	// Create a RateLimiter for each limit
	// Use 'skip' mode for individual limiters since we'll handle blocking here
	for _, limit := range limits {
		limiter, err := NewRateLimiter(limit.Limit, limit.Period, RateLimitSkip)
		if err != nil {
			return nil, err
		}
		mrl.limiters = append(mrl.limiters, limiter)
	}

	return mrl, nil
}

// Acquire attempts to acquire tokens from all rate limiters.
func (mrl *MultiRateLimiter) Acquire(ctx context.Context, tokens int64, onLimitExceeded *RateLimitBehavior) (bool, error) {
	behavior := mrl.onLimitExceeded
	if onLimitExceeded != nil {
		behavior = *onLimitExceeded
	}

	mrl.mu.Lock()
	defer mrl.mu.Unlock()

	// Check if all limiters can satisfy the request
	canProceed := true
	for _, limiter := range mrl.limiters {
		if limiter.GetAvailableTokens() < float64(tokens) {
			canProceed = false
			break
		}
	}

	if canProceed {
		// Acquire from all limiters
		for _, limiter := range mrl.limiters {
			limiter.TryAcquire(tokens)
		}
		return true, nil
	}

	switch behavior {
	case RateLimitBlock:
		// Calculate maximum wait time needed across all limiters
		var maxWaitTime time.Duration
		for _, limiter := range mrl.limiters {
			if limiter.GetAvailableTokens() < float64(tokens) {
				tokensNeeded := float64(tokens) - limiter.GetAvailableTokens()
				waitTime := time.Duration(tokensNeeded * limiter.period.Seconds() / float64(limiter.limit) * float64(time.Second))
				if waitTime > maxWaitTime {
					maxWaitTime = waitTime
				}
			}
		}

		// Release lock while waiting
		mrl.mu.Unlock()

		// Wait with context support
		timer := time.NewTimer(maxWaitTime)
		select {
		case <-timer.C:
			// Continue after wait
		case <-ctx.Done():
			timer.Stop()
			mrl.mu.Lock()
			return false, ctx.Err()
		}

		// Re-acquire lock
		mrl.mu.Lock()

		// Try again after waiting
		allReady := true
		for _, limiter := range mrl.limiters {
			if limiter.GetAvailableTokens() < float64(tokens) {
				allReady = false
				break
			}
		}

		if allReady {
			for _, limiter := range mrl.limiters {
				limiter.TryAcquire(tokens)
			}
			return true, nil
		}
		return false, nil

	case RateLimitRaise:
		return false, ErrRateLimitExceeded

	default: // RateLimitSkip
		return false, nil
	}
}

// TryAcquire attempts to acquire tokens without blocking.
func (mrl *MultiRateLimiter) TryAcquire(tokens int64) bool {
	skip := RateLimitSkip
	acquired, _ := mrl.Acquire(context.Background(), tokens, &skip)
	return acquired
}

// Reset resets all rate limiters to their initial state.
func (mrl *MultiRateLimiter) Reset() {
	mrl.mu.Lock()
	defer mrl.mu.Unlock()

	for _, limiter := range mrl.limiters {
		limiter.Reset()
	}
}

// GetStatus returns the status of all rate limiters.
//
// Returns a slice of (limit, period, availableTokens) for each limiter.
type LimiterStatus struct {
	Limit           int64
	Period          time.Duration
	AvailableTokens float64
}

func (mrl *MultiRateLimiter) GetStatus() []LimiterStatus {
	mrl.mu.RLock()
	defer mrl.mu.RUnlock()

	status := make([]LimiterStatus, len(mrl.limiters))
	for i, limiter := range mrl.limiters {
		status[i] = LimiterStatus{
			Limit:           limiter.limit,
			Period:          limiter.period,
			AvailableTokens: limiter.GetAvailableTokens(),
		}
	}
	return status
}

// TimeUntilReady calculates time until all limiters are ready.
func (mrl *MultiRateLimiter) TimeUntilReady() time.Duration {
	mrl.mu.RLock()
	defer mrl.mu.RUnlock()

	var maxWait time.Duration
	for _, limiter := range mrl.limiters {
		wait := limiter.TimeUntilNextToken()
		if wait > maxWait {
			maxWait = wait
		}
	}
	return maxWait
}

// Wait blocks until all limiters are ready or context is cancelled.
func (mrl *MultiRateLimiter) Wait(ctx context.Context) error {
	_, err := mrl.Acquire(ctx, 1, nil)
	return err
}

// ============================================================================
// Helper Functions
// ============================================================================

// min returns the minimum of two float64 values
func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
