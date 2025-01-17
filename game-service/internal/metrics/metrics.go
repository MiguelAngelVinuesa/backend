package metrics

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/metrics"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
)

var (
	Metrics      *allMetrics
	PrintMetrics = 5 * time.Minute
	StartTime    = now()
)

type allMetrics struct {
	durations *metrics.Durations
}

var rwMutex sync.Mutex

func (m *allMetrics) AddDuration(t metrics.DurationType, start time.Time) {
	if m == nil {
		return
	}
	rwMutex.Lock()
	m.durations.Add(t, time.Since(start))
	rwMutex.Unlock()
}

func (m *allMetrics) Print() {
	d := Metrics.reset()

	for t, m := range d.GetAll() {
		if m != nil && m.Count() > 0 {
			log.Logger.Info("metrics",
				"type", durationNames[t],
				"count", m.Count(),
				"total", m.Total(),
				"min", m.Min(),
				"max", m.Max(),
				"avg", m.Avg(),
			)
		}
	}

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	log.Logger.Info("memory",
		"system", m.pretty(ms.Sys),
		"heap", m.pretty(ms.HeapAlloc),
		"num GC", ms.NumGC,
		"unique sessions", ResetSessions(),
		"uptime", now().Sub(StartTime).String(),
	)
}

func (m *allMetrics) reset() *metrics.Durations {
	rwMutex.Lock()
	d := m.durations.Clone()
	m.durations.Reset()
	rwMutex.Unlock()
	return d
}

func (m *allMetrics) pretty(b uint64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

func InitMetrics() {
	Metrics = &allMetrics{
		durations: metrics.NewDurations(MaxDuration),
	}

	go func() {
		tick := time.NewTicker(PrintMetrics)
		for {
			select {
			case <-tick.C:
				Metrics.Print()
			}
		}
	}()
}

func now() time.Time {
	return time.Now().UTC().Round(time.Second)
}
