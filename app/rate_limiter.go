package app

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

const reqsPerSecond = 5000
const sleepTime = time.Second / reqsPerSecond

var limiter = make(chan struct{}, reqsPerSecond)

func setUpLimiter() {
	for i := 0; i < reqsPerSecond; i++ {
		limiter <- struct{}{}
	}
	go func() {
		for {
			limiter <- struct{}{}
			time.Sleep(sleepTime)
		}
	}()
}

func rateLimiter() {
	<-limiter
}
func perUserRequestLimiter(uid int, ip string) {
	if uid == 0 {
		defaultLimiter.Request("ip:"+ip, 60)
	} else {
		defaultLimiter.Request("user:"+strconv.Itoa(uid), 2000)
	}
}

var defaultLimiter = &specificRateLimiter{
	Map:   make(map[string]chan struct{}),
	Mutex: &sync.RWMutex{},
}

type specificRateLimiter struct {
	Map   map[string]chan struct{}
	Mutex *sync.RWMutex
}

func (s *specificRateLimiter) Request(u string, perMinute int) {
	s.Mutex.RLock()
	c, exists := s.Map[u]
	s.Mutex.RUnlock()
	if !exists {
		c = makePrefilledChan(perMinute)
		s.Mutex.Lock()
		// Now that we have exclusive read and write-access, we want to
		// make sure we don't overwrite an existing channel. Otherwise,
		// race conditions and panic happen.
		if cNew, exists := s.Map[u]; exists {
			c = cNew
			s.Mutex.Unlock()
		} else {
			s.Map[u] = c
			s.Mutex.Unlock()
			<-c
			go s.filler(u, perMinute)
		}
	}
	<-c
}

func (s *specificRateLimiter) filler(el string, perMinute int) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println(r)
		}
	}()

	s.Mutex.RLock()
	c := s.Map[el]
	s.Mutex.RUnlock()
	for {
		select {
		case c <- struct{}{}:
			time.Sleep(time.Minute / time.Duration(perMinute))
		default: // c is full
			s.Mutex.Lock()
			close(c)
			delete(s.Map, el)
			s.Mutex.Unlock()
			return
		}
	}
}

func makePrefilledChan(l int) chan struct{} {
	c := make(chan struct{}, l)
	for i := 0; i < l; i++ {
		c <- struct{}{}
	}
	return c
}
