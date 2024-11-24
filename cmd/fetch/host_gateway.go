package main

import (
	"sync"
	"time"

	"go.uber.org/zap"
)

type HostGateway struct {
	politeSleep time.Duration
	last        map[string]time.Time
	m           *sync.RWMutex
}

func NewHostGateway(politeSleep time.Duration) *HostGateway {
	return &HostGateway{
		politeSleep: politeSleep,
		last:        make(map[string]time.Time),
		m:           &sync.RWMutex{},
	}
}

// When I see a host, I should check if it has a non-zero counter.
func (hsm *HostGateway) GoodToGo(host string) bool {
	hsm.m.Lock()
	defer hsm.m.Unlock()

	// Lets see if they're in the map
	lastHit, ok := hsm.last[host]

	if ok {
		// We have seen this host before
		timeSinceLastHit := time.Since(lastHit)
		// If it has been more than polite sleep since the last hit,
		// then the caller is good to go.
		isGoodToGo := timeSinceLastHit > hsm.politeSleep

		zap.L().Debug("gateway host is good to go",
			zap.Int64("timeSinceLastHit", int64(timeSinceLastHit)),
			zap.Int64("politeSleep", int64(hsm.politeSleep)),
		)

		if isGoodToGo {
			// If they are good to go, we'll update the map.
			hsm.last[host] = time.Now()
		}
		return isGoodToGo
	} else {
		// We have not seen this host before
		// Therefore, add them to the map, and they're good to go.
		zap.L().Debug("gateway: host never seen before")
		hsm.last[host] = time.Now()
		return true
	}
}

func (hsm *HostGateway) TimeRemaining(host string) time.Duration {
	// This is a read operation.
	hsm.m.RLock()
	defer hsm.m.RUnlock()

	lastHit, ok := hsm.last[host]

	// We must be in the map. But, we should check.
	if ok {
		// If we're in the map, there must be some time until the
		// next possible access to the host.
		until := time.Until(lastHit.Add(hsm.politeSleep))

		zap.L().Debug("gateway time remaining",
			zap.Int64("until", int64(until)))

		if until > 0 {
			return until
		} else {
			return time.Duration(0 * time.Millisecond)
		}
	} else {
		// If someone asks for a host that is not in the map, we'll tell them
		// there are 0 milliseconds until the host is ready.
		return time.Duration(0 * time.Millisecond)
	}
}
