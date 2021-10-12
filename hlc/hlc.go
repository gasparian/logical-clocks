package hlc

import (
	"sync"
	"time"
)

func getWallTimeMs() int64 {
	return time.Now().UnixNano() / 1e6
}

type Hybrid struct {
	mu              sync.RWMutex
	WallClockTimeMs int64
	Ticks           int64
}

func New(systemTime, ticks int64) *Hybrid {
	return &Hybrid{
		WallClockTimeMs: systemTime,
		Ticks:           ticks,
	}
}

func NewNow(ticks int64) *Hybrid {
	return &Hybrid{
		WallClockTimeMs: getWallTimeMs(),
		Ticks:           ticks,
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
	currentTime := getWallTimeMs()
	// NOTE: check in case clock goes backwards fue to synchronization
	if h.WallClockTimeMs >= currentTime {
		h.addTicks(1)
	} else {
		h.WallClockTimeMs = currentTime
		h.Ticks = 0
	}
	return Hybrid{
		WallClockTimeMs: h.WallClockTimeMs,
		Ticks:           h.Ticks,
	}
}

func Compare(a, b *Hybrid) int {
	if (a.WallClockTimeMs == b.WallClockTimeMs) && (a.Ticks == b.Ticks) {
		return 0
	}
	if (a.WallClockTimeMs == b.WallClockTimeMs && a.Ticks > b.Ticks) ||
		(a.WallClockTimeMs > b.WallClockTimeMs) {
		return 1
	}
	return -1
}

func max(times ...*Hybrid) *Hybrid {
	maxTime := times[0]
	for _, time := range times {
		cmp := Compare(time, maxTime)
		if cmp == 1 {
			maxTime = time
		}
	}
	return maxTime
}

func (h *Hybrid) Tick(requestTime *Hybrid) {
	h.mu.Lock()
	defer h.mu.Unlock()
	hybridNow := NewNow(-1)
	latestTime := max(hybridNow, requestTime, h)
	latestTime.addTicks(1)
	h.WallClockTimeMs = latestTime.WallClockTimeMs
	h.Ticks = latestTime.Ticks
}

func (h *Hybrid) GetCompactTimestampMs() int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return (h.WallClockTimeMs >> 16 << 16) | (h.Ticks << 48 >> 48)
}
