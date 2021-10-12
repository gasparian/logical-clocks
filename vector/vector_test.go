package vector

import (
	"sync"
	"testing"
)

func TestCompare(t *testing.T) {
	_, c := New(0, 3)
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
	_, c1 := New(1, 3)
	c1.AddTicks(1)
	cmp = Compare(c, c1)
	if cmp != 0 {
		t.Fatalf("Timestamps must be concurrent: %v != %v", c, c1)
	}
}

func TestConcurrentAdd(t *testing.T) {
	_, c := New(0, 3)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.AddTicks(1)
		}()
	}
	wg.Wait()
	now := c.Now()
	if now.GetTicks() < 10 {
		t.Fatal("There must be exactly 10 ticks")
	}
}

func TestConcurrentTick(t *testing.T) {
	_, c := New(0, 3)
	_, req := New(1, 3)
	req.AddTicks(1)
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			c.Tick(req)
		}()
	}
	wg.Wait()
	now := c.Now()
	if now.Ticks[0] != 10 && now.Ticks[1] != 1 {
		t.Fatal("Vector must be equal to {10, 1}")
	}
}
