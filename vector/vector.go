package vector

import (
	"errors"
	"sync"
)

type Vector struct {
	mu    sync.RWMutex
	Self  int
	Ticks []int64
}

func New(self int, nNodes int) (error, *Vector) {
	if self < 0 || self >= nNodes {
		return errors.New("Self node id must be > 0 and < total number of nodes"), nil
	}
	return nil, &Vector{
		Self:  self,
		Ticks: make([]int64, nNodes),
	}
}

func (v *Vector) addTicks(ticks int64) {
	v.Ticks[v.Self] += ticks
}

func (v *Vector) AddTicks(ticks int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.Ticks[v.Self] += ticks
}

func (v *Vector) GetTicks() int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return v.Ticks[v.Self]
}

func (v *Vector) Now() Vector {
	cpy := make([]int64, len(v.Ticks))
	copy(cpy, v.Ticks)
	return Vector{
		Self:  v.Self,
		Ticks: cpy,
	}
}

func (v *Vector) Tick(requestTime *Vector) error {
	v.mu.Lock()
	defer v.mu.Unlock()
	if len(v.Ticks) != len(requestTime.Ticks) {
		return errors.New("Vectors lengths must be equal")
	}
	for i := range v.Ticks {
		if requestTime.Ticks[i] > v.Ticks[i] {
			v.Ticks[i] = requestTime.Ticks[i]
		}
	}
	v.addTicks(1)
	return nil
}

func Compare(a, b *Vector) int {
	equal, less := 0, 0
	for i := range a.Ticks {
		if a.Ticks[i] == b.Ticks[i] {
			equal++
		} else if a.Ticks[i] < b.Ticks[i] {
			less++
		}
	}
	if len(a.Ticks) == equal {
		return 0
	}
	if len(a.Ticks) == (equal + less) {
		return -1
	}
	return 1
}
