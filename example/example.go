package main

import (
	"runtime"
	"sync"
	"time"

	"github.com/zweite/rate"
)

func main() {
	runtime.GOMAXPROCS(1)
	r := rate.NewRate(100)
	go r.Run()

	go func() {
		qpsChan := r.QPS()
		for {
			select {
			case qps := <-qpsChan:
				println("QPS", qps)
			}
		}
	}()

	wg := sync.WaitGroup{}
	count := 0
	num := 100000

	for i := 0; i < num; i++ {
		time.Sleep(5 * time.Millisecond)
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			if r.GetToken(1 * time.Millisecond) {
				count++
			}
		}(i)
	}
	wg.Wait()

	println("count", count)
}
