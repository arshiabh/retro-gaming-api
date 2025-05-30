package ratelimiter

import "time"

type FixedWindowLimiter struct {
	RequestPerTime int
	TimeFrame      time.Duration
}

func NewFixedWindowLimiter(request int, timeFrame time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		RequestPerTime: request,
		TimeFrame:      timeFrame,
	}
}

func (r *FixedWindowLimiter) Allow(ip string) (bool, time.Duration)
