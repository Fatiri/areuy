package times

import "time"

type Time interface {
	Now(timeGMT *int) time.Time
	TimeStampToDateStr(timeStr, layout string) string
	TimeStampToDate(timeStr, layout string) time.Time
}
