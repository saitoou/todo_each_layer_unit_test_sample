package utils

import "time"

func JstLocation() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
