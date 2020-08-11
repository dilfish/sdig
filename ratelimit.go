// sean at shanghai

package main

import (
	"sync"
	"time"
)

type QPSRateLimiter struct {
	count int
	init  int
	lock  sync.Mutex
	start time.Time
}

func NewQPSRateLimiter(c int) *QPSRateLimiter {
	if c <= 1 {
		c = 1
	}
	return &QPSRateLimiter{count: c, init: c, start: time.Now()}
}

func (q *QPSRateLimiter) GetToken() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	now := time.Now()
	if now.Sub(q.start) > time.Second {
		q.start = now
		q.count = q.init
	}
	q.count = q.count - 1
	if q.count <= 0 {
		return false
	}
	return true
}

func (q *QPSRateLimiter) WaitForToken() {
	for {
		if q.GetToken() == false {
			time.Sleep(time.Millisecond * 100)
		} else {
			return
		}
	}
}
