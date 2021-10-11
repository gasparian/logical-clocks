package hlc

import (
	"sync"
	"time"
)

type Hybrid struct {
	mu    sync.Mutex
	Time  int64 // Usually time in ms
	Ticks int64
}

func New(systemTime, ticks int64) *Hybrid {
	return &Hybrid{
		Time:  systemTime,
		Ticks: ticks,
	}
}

func (h *Hybrid) addTicks(ticks int64) {
	h.Ticks += ticks
}

func (h *Hybrid) AddTicks(ticks int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Ticks += ticks
}

func (h *Hybrid) Now() Hybrid {
	h.mu.Lock()
	defer h.mu.Unlock()
	currentTime := time.Now().UnixNano() / 1000
	if h.Time >= currentTime {
		h.addTicks(1)
	} else {
		h.Time = currentTime
		h.Ticks = 0
	}
	return Hybrid{
		Time:  h.Time,
		Ticks: h.Ticks,
	}
}

func Compare(a, b *Hybrid) int {
	if (a.Time == b.Time) && (a.Ticks == b.Ticks) {
		return 0
	}
	if (a.Time == b.Time && a.Ticks > b.Ticks) ||
		(a.Time > b.Time) {
		return 1
	}
	return -1
}

func max(times ...*Hybrid) *Hybrid {
	maxTime := times[0]
	for _, time := range times {
		cmp := Compare(maxTime, time)
		if cmp == 1 {
			maxTime = time
		}
	}
	return maxTime
}

func (h *Hybrid) Tick(requestTime *Hybrid) {
	h.mu.Lock()
	defer h.mu.Unlock()
	current := time.Now().UnixNano() / 1000
	hybridNow := New(current, -1)
	latestTime := max(hybridNow, requestTime, h)
	latestTime.addTicks(1)
	h.Time = latestTime.Time
	h.Ticks = latestTime.Ticks
}
