package usage

import (
	"time"

	"go.uber.org/atomic"
)

type (
	Tracker interface {
		Utilization() float64
	}
	backgroundTracker struct {
		tracker     Tracker
		utilisation *atomic.Float64
	}
)

func NewBackgroundTracker(t Tracker, refreshInterval time.Duration) Tracker {
	b := &backgroundTracker{
		tracker:     t,
		utilisation: atomic.NewFloat64(0),
	}
	go b.runRefresher(refreshInterval)
	return b
}

func (b *backgroundTracker) Utilization() float64 {
	return b.utilisation.Load()
}

func (b *backgroundTracker) runRefresher(interval time.Duration) {
	tick := time.NewTicker(interval)
	tickCh := tick.C
	defer tick.Stop()
	for {
		b.refreshUtilization()
		<-tickCh
	}
}

func (b *backgroundTracker) refreshUtilization() {
	util := b.tracker.Utilization()
	b.utilisation.Store(util)
}
