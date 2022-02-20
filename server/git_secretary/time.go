package git_secretary

import "time"

func GetDateHash(year int, month time.Month, day int) int {
	return year*10000 + int(month)*100 + day
}
