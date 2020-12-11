package counter

import (
	"fmt"
	"github.com/suikammd/rate-limiter/req"
	"math/rand"
	"sync"
	"time"
)

type Node struct {
	next    *Node
	counter int

	start time.Time
	end   time.Time
}

type LinkedList struct {
	Size int
	Head *Node
	Tail *Node
}

type SlidingCounter struct {
	Limit    int
	interval int
	list     *LinkedList

	Ticker *time.Ticker
	stopCh chan bool

	sync.Mutex
}

func NewSlidingCounter(limit, slot int) SlidingCounter {
	if limit < 0 || slot < 0 {
		panic("illegal params")
	}

	interval := 1000 / slot
	start := time.Now()
	cur := &Node{
		counter: 0,
		start:   start,
		end:     start.Add(time.Duration(interval) * time.Millisecond),
	}
	head := cur
	start = cur.end
	for i := 0; i < slot-1; i++ {
		end := start.Add(time.Duration(interval) * time.Millisecond)
		cur.next = &Node{
			counter: 0,
			start:   start,
			end:     end,
		}
		cur = cur.next
		start = end
	}
	tail := cur

	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	return SlidingCounter{
		Limit:    limit,
		Ticker:   ticker,
		interval: interval,
		stopCh:   make(chan bool),
		list: &LinkedList{
			Size: slot,
			Head: head,
			Tail: tail,
		},
	}
}

func (s *SlidingCounter) resetCounter() {
	for {
		select {
		case <-s.Ticker.C:
			s.Lock()
			s.list.removeHead()
			s.list.addToTail(s.interval)
			s.Unlock()
		case <-s.stopCh:
			return
		}
	}
}

func (s *SlidingCounter) ServeRequest(r req.Request) bool {
	s.Lock()
	defer s.Unlock()

	cnt := s.calcReqCnt()
	if cnt > s.Limit {
		r.Func(false)
		return false
	}

	idx := int((time.Now().UnixNano() - s.list.Head.start.UnixNano()) / (1000 * 1000 * int64(s.interval)))
	s.list.index(idx).counter++
	r.Func(true)
	return true
}

func (s SlidingCounter) Start() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		s.resetCounter()
	}()

	go func() {
		cnt := 0
		for i := 0; i < 10; i++ {
			reqs := req.GenReq(rand.Intn(5))
			for _, req := range reqs {
				fmt.Println(cnt)
				s.ServeRequest(req)
				cnt++
			}
			time.Sleep(time.Duration(rand.Intn(50)*100) * time.Millisecond)
		}
	}()

	time.Sleep(5 * time.Second)
	s.stopCh <- true
	wg.Wait()
}

func (s *SlidingCounter) calcReqCnt() int {
	cnt := 0
	cur := s.list.Head
	for i := 0; i < s.list.Size-1; i++ {
		cnt = cnt + cur.counter
		cur = cur.next
	}
	return cnt
}

func (l *LinkedList) removeHead() {
	l.Head = l.Head.next
}

func (l *LinkedList) addToTail(interval int) {
	l.Tail.next = &Node{
		counter: 0,
		start:   l.Tail.end,
		end:     l.Tail.end.Add(time.Duration(interval) * time.Millisecond),
	}
	l.Tail = l.Tail.next
}

func (l *LinkedList) index(idx int) *Node {
	if idx < 0 || idx >= l.Size {
		panic("illegal index")
	}

	cur := l.Head
	for i := 0; i < idx; i++ {
		cur = cur.next
	}
	return cur
}
