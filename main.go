package main

import "github.com/suikammd/rate-limiter/counter"

func main()  {
	//counter.NewFixedWindowCounter(1, 100).Start()
	counter.NewSlidingCounter(10, 10).Start()
}
