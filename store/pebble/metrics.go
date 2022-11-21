package pebble

import (
	"context"
	"time"
)

const (
	defaultMetricsInterval = 5 * time.Second
)

func (s *PebbleStore) RunMetrics(ctx context.Context, interval time.Duration) {
	if interval == 0 {
		interval = defaultMetricsInterval
	}

	tick := time.NewTicker(interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.emitMetrics()
		}
	}
}

func (s *PebbleStore) emitMetrics() {
	//stats := s.db.Metrics()

	//// freelist metrics
	//metrics.SetGauge([]string{"raft", "boltdb", "numFreePages"}, float32(newStats.FreePageN))
	//metrics.SetGauge([]string{"raft", "boltdb", "numPendingPages"}, float32(newStats.PendingPageN))
	//metrics.SetGauge([]string{"raft", "boltdb", "freePageBytes"}, float32(newStats.FreeAlloc))
	//metrics.SetGauge([]string{"raft", "boltdb", "freelistBytes"}, float32(newStats.FreelistInuse))
	//
	//// txn metrics
	//metrics.IncrCounter([]string{"raft", "boltdb", "totalReadTxn"}, float32(stats.TxN))
	//metrics.SetGauge([]string{"raft", "boltdb", "openReadTxn"}, float32(newStats.OpenTxN))
	//
	//// tx stats
	//metrics.SetGauge([]string{"raft", "boltdb", "txstats", "pageCount"}, float32(newStats.TxStats.PageCount))
	//metrics.SetGauge([]string{"raft", "boltdb", "txstats", "pageAlloc"}, float32(newStats.TxStats.PageAlloc))
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "cursorCount"}, float32(stats.TxStats.CursorCount))
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "nodeCount"}, float32(stats.TxStats.NodeCount))
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "nodeDeref"}, float32(stats.TxStats.NodeDeref))
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "rebalance"}, float32(stats.TxStats.Rebalance))
	//metrics.AddSample([]string{"raft", "boltdb", "txstats", "rebalanceTime"}, float32(stats.TxStats.RebalanceTime.Nanoseconds())/1000000)
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "split"}, float32(stats.TxStats.Split))
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "spill"}, float32(stats.TxStats.Spill))
	//metrics.AddSample([]string{"raft", "boltdb", "txstats", "spillTime"}, float32(stats.TxStats.SpillTime.Nanoseconds())/1000000)
	//metrics.IncrCounter([]string{"raft", "boltdb", "txstats", "write"}, float32(stats.TxStats.Write))
	//metrics.AddSample([]string{"raft", "boltdb", "txstats", "writeTime"}, float32(stats.TxStats.WriteTime.Nanoseconds())/1000000)
	//return &newStats
}
