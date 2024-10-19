package cpu

import (
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/fhke/loadsheding-go/usage"
	"github.com/tklauser/go-sysconf"
)

type timeTracker struct {
	clockTicksCnf float64

	lastCall       time.Time
	lastUsageTicks float64
}

func New() (usage.Tracker, error) {
	t := &timeTracker{}
	if err := t.loadConf(); err != nil {
		return nil, fmt.Errorf("error loading config: %w", err)
	}

	return t, nil
}

func (t *timeTracker) Utilization() float64 {
	// get CPU ticks
	tms := new(syscall.Tms)
	_, err := syscall.Times(tms)
	if err != nil {
		log.Printf("[ERROR] Error gathering CPU times: %s", err.Error())
		return 0
	}
	cumulativeTicks := float64(tms.Stime + tms.Utime)
	elapsedTicks := cumulativeTicks - t.lastUsageTicks
	t.lastUsageTicks = cumulativeTicks

	// get elapsed wall time
	timeNow := time.Now()
	elapsedWall := timeNow.Sub(t.lastCall)
	t.lastCall = timeNow

	elapsedCPUTime := elapsedTicks / (elapsedWall.Seconds() * t.clockTicksCnf)
	return elapsedCPUTime
}

func (t *timeTracker) loadConf() error {
	ticks, err := sysconf.Sysconf(sysconf.SC_CLK_TCK)
	if err != nil {
		return fmt.Errorf("error loading sysconf: %w", err)
	}
	t.clockTicksCnf = float64(ticks)
	return nil
}
