package util

import "time"

// MonotonicMs Gets the current monotonic time in milliseconds.
func MonotonicMs() int64 {
	return time.Now().UnixNano() / 1e6
}

func VerifyGroupId(groupId string) error {
	// TODO: to implement
	return nil
}
