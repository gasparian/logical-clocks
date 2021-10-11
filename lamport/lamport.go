package lamport

import (
	"sync"
)

type Lamport struct {
	mu    sync.Mutex
	Name  string
	Ticks int64
}

func New(name string, ticks int64) *Lamport {
	return &Lamport{
		Name:  name,
		Ticks: ticks,
	}
}

func (lt *Lamport) addTicks(ticks int64) {
	lt.Ticks += ticks
}

func (lt *Lamport) AddTicks(ticks int64) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.Ticks += ticks
}

func (lt *Lamport) Now() Lamport {
	return Lamport{
		Name:  lt.Name,
		Ticks: lt.Ticks,
	}
}

func Compare(a, b *Lamport) int {
	if (a.Ticks == b.Ticks) && (a.Name == b.Name) {
		return 0
	}
	if (a.Ticks > b.Ticks) ||
		((a.Ticks == b.Ticks) && (a.Name > b.Name)) {
		return 1
	}
	return -1
}

func (lt *Lamport) Tick(requestTime *Lamport) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	if cmp := Compare(requestTime, lt); cmp == 1 {
		lt.Ticks = requestTime.Ticks
	}
	lt.addTicks(1)
}
