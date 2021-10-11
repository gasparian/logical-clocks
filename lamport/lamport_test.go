package lamport

import (
	"sync"
	"testing"
)

func TestCompare(t *testing.T) {
	lt := New("a", 0)
	now := lt.Now()
	cmp := Compare(lt, &now)
	if cmp != 0 {
		t.Fatal("Timestamps should be equal")
	}
	lt.AddTicks(1)
	cmp = Compare(lt, &now)
	if cmp < 1 {
		t.Fatal("Original timestamp now must be higher than it's copy")
	}
}

func TestConcurrentAdd(t *testing.T) {
	lt := New("a", 0)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			lt.AddTicks(1)
		}()
	}
	wg.Wait()
	now := lt.Now()
	if now.Ticks < 10 {
		t.Fatal("There must be exactly 10 ticks")
	}
}

func TestConcurrentTick(t *testing.T) {
	lt := New("a", 0)
	req := New("b", 0)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			lt.Tick(req)
		}()
	}
	wg.Wait()
	now := lt.Now()
	if now.Ticks < 10 {
		t.Fatal("There must be exactly 10 ticks")
	}
}
