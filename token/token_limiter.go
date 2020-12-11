package token

import (
	"fmt"
	"github.com/suikammd/rate-limiter/req"
	"time"
)

type Bucket struct {
	limit  int
	bucket chan bool

	ticker *time.Ticker
	stopCh chan bool
}

func NewTokenBucket(limit int) Bucket {
	ticker := time.NewTicker(time.Duration(1000/limit) * time.Millisecond)
	return Bucket{
		limit:  limit,
		bucket: make(chan bool, limit),
		ticker: ticker,
		stopCh: make(chan bool),
	}
}

func (b *Bucket) Start() {
	go b.addToken()
	time.Sleep(time.Second)
	fmt.Println(time.Now())

	go func() {
		cnt := 0
		for i := 0; i < 10; i++ {
			reqs := req.GenReq(100)
			for _, req := range reqs {
				select {
				case <-b.bucket:
					req.Func(true)
					cnt++
				}
			}
		}
	}()

	time.Sleep(5 * time.Second)
	close(b.stopCh)
}

func (b *Bucket) addToken() {
	for {
		select {
		case <-b.ticker.C:
			b.bucket <- true
		case <-b.stopCh:
			return
		}
	}
}
