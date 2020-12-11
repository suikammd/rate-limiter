package counter

import (
	"github.com/suikammd/rate-limiter/req"
	"sync"
	"time"
)

type FixedCounter struct {
	Counter int
	Limit   int
	Ticker  *time.Ticker

	StopCh chan bool
	sync.Mutex
}

func NewFixedWindowCounter(window, limit int) FixedCounter {
	if window < 0 || limit < 0 {
		panic("illegal params")
	}

	ticker := time.NewTicker(time.Duration(window) * time.Second)
	return FixedCounter{
		Counter: 0,
		Limit:   limit,
		Ticker:  ticker,

		StopCh: make(chan bool),
	}
}

func (f *FixedCounter) resetCounter() {
	for {
		select {
		case <-f.Ticker.C:
			f.Lock()
			f.Counter = 0
			f.Unlock()
		case <-f.StopCh:
			return
		}
	}
}

func (f *FixedCounter) ServeRequest(r req.Request) bool {
	f.Lock()
	defer f.Unlock()

	if f.Counter == f.Limit {
		r.Func(false)
		return false
	}

	f.Counter++
	r.Func(true)
	return true
}

func (f FixedCounter) Start() {
	go func() {
		f.resetCounter()
	}()

	reqs := req.GenReq(101)
	for i := 0; i < len(reqs); i++ {
		go f.ServeRequest(reqs[i])
	}
	time.Sleep(2 * time.Second)
	for i := 0; i < len(reqs); i++ {
		go f.ServeRequest(reqs[i])
	}
	f.StopCh <- true
	time.Sleep(2 * time.Second)
}
