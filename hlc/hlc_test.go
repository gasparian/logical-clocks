package hlc

import (
	"sync"
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	c := NewNow(0)
	now := c.Now()
	cmp := Compare(c, &now)
	if cmp != 0 {
		t.Fatal("Timestamps should be equal")
	}
	c.AddTicks(1)
	cmp = Compare(c, &now)
	if cmp < 1 {
		t.Fatal("Original timestamp now must be higher than it's copy")
	}
	time.Sleep(1)
	afterSleep := c.Now()
	cmp = Compare(&now, &afterSleep)
	if cmp > -1 {
		t.Fatal("Time with greater wall clock time should be larger then original timestamp")
	}
}

func TestConcurrentAdd(t *testing.T) {
	currentWallTime := getWallTimeMs()
	c1 := New(currentWallTime, 0)
	c2 := New(currentWallTime, 1)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c1.AddTicks(1)
		}()
	}
	wg.Wait()
	cmp := Compare(c1, c2)
	if cmp < 1 {
		t.Fatal("Timestamp with more ticks must be greater")
	}
}

func TestConcurrentTick(t *testing.T) {
	c := NewNow(0)
	req := New(getWallTimeMs()+1000, 10)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.Tick(req)
		}()
	}
	wg.Wait()
	if c.Ticks < 20 {
		t.Fatal("There must be exactly 20 ticks")
	}
}
