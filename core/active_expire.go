package core

import (
	"log"
	"time"
)

var activeExpirySampleSize int = 20

// Deletes all the expired keys - active way
// Sampling approach: https://redis.io/commands/expire/
func ScanAndDeleteExpiredKeys() {
	for {
		frac := scanAndDeleteKeys()
		// Break if the sample had less than 25% keys expired
		if frac < 0.25 {
			break
		}
	}
	log.Println("deleted the expired but undeleted keys. total keys", len(store.store))
}

func scanAndDeleteKeys() float32 {
	sampleSize := activeExpirySampleSize
	expiredCount := 0

	for k, v := range store.store {
		if v.expiresAt == -1 {
			continue
		}
		sampleSize--
		if v.expiresAt <= time.Now().UnixMilli() {
			delete(store.store, k)
			expiredCount++
		}
		if sampleSize == 0 {
			break
		}
	}

	return float32(expiredCount / activeExpirySampleSize)
}