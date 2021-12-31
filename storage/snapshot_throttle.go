package storage

type SnapshotThrottle interface {
	ThrottledByThroughput(bytes int64) int64
}
