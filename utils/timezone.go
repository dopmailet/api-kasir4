package utils

import "time"

// GetTimezone mengambil timezone dari query parameter dengan validasi.
// Default: Asia/Makassar
func GetTimezone(tz string) string {
	if tz == "" {
		return "Asia/Makassar"
	}
	_, err := time.LoadLocation(tz)
	if err != nil {
		return "Asia/Makassar"
	}
	return tz
}
