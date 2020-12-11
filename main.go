package main

import (
	"github.com/suikammd/rate-limiter/token"
)

func main()  {
	//counter.NewFixedWindowCounter(1, 100).Start()
	//counter.NewSlidingCounter(10, 10).Start()
	b := token.NewTokenBucket(1000)
	b.Start()
}
