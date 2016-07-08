package main

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/zweite/rate"
)

func main() {
	runtime.GOMAXPROCS(1)
	// set the QPS
	r := rate.NewRate(200)
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

	seed := 5 // if seed down then QPS will up
	for i := 0; i < num; i++ {
		time.Sleep(time.Duration(rand.Intn(seed)) * time.Millisecond)
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			if r.GetToken() {
				count++
			}
		}(i)
	}
	wg.Wait()

	println("count", count)
}
