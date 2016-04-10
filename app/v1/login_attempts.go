package v1

import (
	"sync"
	"time"
)

type failedAttempt struct {
	attempt time.Time
	ID      int
}

var failedAttempts []failedAttempt
var failedAttemptsMutex = new(sync.RWMutex)

// removeUseless removes the expired attempts in failedAttempts
func removeUseless() {
	for {
		failedAttemptsMutex.RLock()
		var localCopy = make([]failedAttempt, len(failedAttempts))
		copy(localCopy, failedAttempts)
		failedAttemptsMutex.RUnlock()
		var newStartFrom int
		for k, v := range localCopy {
			if time.Since(v.attempt) > time.Minute*10 {
				newStartFrom = k + 1
			} else {
				break
			}
		}
		copySl := localCopy[newStartFrom:]
		failedAttemptsMutex.Lock()
		failedAttempts = make([]failedAttempt, len(copySl))
		for i, v := range copySl {
			failedAttempts[i] = v
		}
		failedAttemptsMutex.Unlock()
		time.Sleep(time.Minute * 10)
	}
}

func addFailedAttempt(uid int) {
	failedAttemptsMutex.Lock()
	failedAttempts = append(failedAttempts, failedAttempt{
		attempt: time.Now(),
		ID:      uid,
	})
	failedAttemptsMutex.Unlock()
}

func nFailedAttempts(uid int) int {
	var count int
	failedAttemptsMutex.RLock()
	for _, i := range failedAttempts {
		if i.ID == uid && time.Since(i.attempt) < time.Minute*10 {
			count++
		}
	}
	failedAttemptsMutex.RUnlock()
	return count
}
