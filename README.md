![main tests](https://github.com/gasparian/logical-clocks/actions/workflows/test.yml/badge.svg?branch=main)  

# logical-clocks  
In the distributed systems there is a problem of defining casual order of events. Usually you need that to implement MVCC or CRDTs. And you can't always rely only on physical clocks since it is very hard to achieve tight clocks synchronization in real world.  
But there is a very "elegant" solution to that: use logical clocks, where we can also rely on counting events occuring in the system.  
So here you will find my naive go implementation of three basic logical clocks types: [Lamport clock](https://lamport.azurewebsites.net/pubs/time-clocks.pdf), [Vector clock](https://fileadmin.cs.lth.se/cs/Personal/Amr_Ergawy/dist-algos-papers/4.pdf) and [Hybrid logical clock](https://cse.buffalo.edu/tech-reports/2014-04.pdf).  

## Usage  
Run `make init` to replace pre-commit hook into `.git` folder.  
To test the clocks implementations run `make test`.  
Getting the package:  
```
go get github.com/gasparian/logical-clocks
```  

Then you can start using clocks:  
```go
import (
    hlc "github.com/gasparian/logical-clocks/hlc"
)

// Instantiate hybrid clock with the current wall clock timestamp
hc := hlc.NewNow(0)
// ...handle request from other service
req := ...
// Apply tick using the hybrid timestamp that came from another service
hc.Tick(req.Time)
// Get current time
currentTime := hc.Now()
// Compare two timestamps if needed
cmp := hlc.Compare(currentTime, req.Time)
// ...
```  

### API  
All clocks shares the same API:  
  - `New() *SomeClock` creates a new clock object;  
  - `clock.Now() SomeClock` returns a copy of the current state of the clock;  
  - `clock.AddTicks(ticks int)` adds the specified number to the internal counter (thread-safe);  
  - `clock.Tick(requestTime *SomeClock)` moves clock forward, comparing to the incoming value;  
  - `Compare(a, b *SomeClock)` compares the current states of two clocks, returns -1, 0 or 1;  

### References  
 - [Martin Fowler's blog](https://martinfowler.com/articles/patterns-of-distributed-systems/)  
