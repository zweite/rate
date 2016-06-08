package rate

import "time"

// Rate control speed for QPS
type Rate struct {
	tokenBucket chan struct{}
	stop        chan struct{}
	restart     chan struct{}
	limit       int
	qps         int64
}

// NewRate Initialization entry
func NewRate(limit int) *Rate {
	return &Rate{
		tokenBucket: make(chan struct{}, 1),
		stop:        make(chan struct{}),
		restart:     make(chan struct{}),
		limit:       limit,
	}
}

// Run new token to bucket
func (r *Rate) Run() {
	tick := time.NewTicker(time.Second / time.Duration(r.limit))
	for {
		select {
		case <-tick.C:
			r.tokenBucket <- struct{}{}
		case <-r.stop:
			r.restart <- struct{}{}
			return
		}
	}
}

// Stop stop new token to bucket
func (r *Rate) Stop() {
	r.stop <- struct{}{}
}

// Restart for new limit
func (r *Rate) Restart(limit int) {
	r.limit = limit
	r.Stop()
	<-r.restart
	go r.Run()
}

// GetToken control QPS
func (r *Rate) GetToken(exp time.Duration) bool {
	timer := time.NewTimer(exp)
	select {
	case <-r.tokenBucket:
		r.qps++
		return true
	case <-timer.C:
		return false
	}
}

// QPS get rate QPS
func (r *Rate) QPS() <-chan int64 {
	tick := time.NewTicker(1 * time.Second)
	qps := make(chan int64)
	var zero int64

	go func() {
		for {
			select {
			case <-tick.C:
				qps <- r.qps
				r.qps = zero
			}
		}
	}()
	return qps
}
