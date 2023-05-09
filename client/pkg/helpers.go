package pkg

import "time"

// MakeTimestamp returns timestamp
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
