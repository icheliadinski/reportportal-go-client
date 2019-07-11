package rp

import (
	"time"
)

func toTimestamp(t time.Time) int64 {
	return t.Unix() * int64(time.Microsecond)
}
