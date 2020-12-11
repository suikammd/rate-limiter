package req

import "fmt"

type Request struct {
	Func func(ac bool)
}

func GenReq(total int) []Request {
	reqs := make([]Request, total)
	for i := 0; i < total; i++ {
		j := i
		reqs[i] = Request{Func: func(ac bool) {
			if ac {
				fmt.Printf("%d req accecpted\n", j)
			} else {
				fmt.Printf("%d req rejected\n", j)
			}
		}}
	}
	return reqs
}
