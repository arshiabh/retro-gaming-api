package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	clients map[string]int
	sync.RWMutex
	window time.Duration
	limit  int
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		window:  window,
		limit:   limit,
		clients: make(map[string]int),
	}
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exist := rl.clients[ip]
	if !exist || count < rl.limit {
		rl.Lock()
		if !exist {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}

	return false, rl.window
}

func (rl *FixedWindowLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}
