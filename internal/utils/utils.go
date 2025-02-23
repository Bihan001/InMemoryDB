package utils

import "time"

func GetCurrentTime() uint64 {
	return uint64(time.Now().UnixMilli())
}
