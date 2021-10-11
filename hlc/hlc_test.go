package hlc

import (
	"time"
	// "sync"
	"testing"
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

// func TestConcurrentAdd(t *testing.T) {
// 	lt := New("a", 0)
// 	wg := sync.WaitGroup{}
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			defer wg.Done()
// 			lt.AddTicks(1)
// 		}()
// 	}
// 	wg.Wait()
// 	now := lt.Now()
// 	if now.Ticks < 10 {
// 		t.Fatal("There must be exactly 10 ticks")
// 	}
// }

// func TestConcurrentTick(t *testing.T) {
// 	lt := New("a", 0)
// 	req := New("b", 0)
// 	wg := sync.WaitGroup{}
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go func() {
// 			defer wg.Done()
// 			lt.Tick(req)
// 		}()
// 	}
// 	wg.Wait()
// 	now := lt.Now()
// 	if now.Ticks < 10 {
// 		t.Fatal("There must be exactly 10 ticks")
// 	}
// }
