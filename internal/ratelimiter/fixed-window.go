package ratelimiter

import (
	"context"
	"log"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	clients map[string]clientInfo
	sync.RWMutex
	window time.Duration
	limit  int
}

type clientInfo struct {
	firstSeen time.Time
	count     int
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		window:  window,
		limit:   limit,
		clients: make(map[string]clientInfo),
	}
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	rl.Lock()
	defer rl.Unlock()

	client, exists := rl.clients[ip]
	now := time.Now()

	if !exists || now.Sub(client.firstSeen) > rl.window {
		rl.clients[ip] = clientInfo{count: 1, firstSeen: now}
		return true, 0
	}

	if client.count < rl.limit {
		client.count++
		// update client
		rl.clients[ip] = client
		return true, 0
	}

	remaining := rl.window - now.Sub(client.firstSeen)
	return false, remaining

}

func (rl *FixedWindowLimiter) StartCleanup(ctx context.Context, wg *sync.WaitGroup, interval time.Duration) {
	defer wg.Done()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("closing rate limiter cleanup")
			return
		case <-ticker.C:
			rl.Lock()
			now := time.Now()
			for ip, client := range rl.clients {
				if now.Sub(client.firstSeen) > rl.window {
					delete(rl.clients, ip)
				}
			}
			rl.Unlock()
		}
	}

}
