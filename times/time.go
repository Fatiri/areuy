package times

import (
	"strconv"
	"time"
)

type timesCustomImpl struct {
	gmt int
}

func ProvideNewTimesCustom() Time {
	// default GMT +7 (Asia/Jakarta)
	return &timesCustomImpl{gmt: 7}
}

func (t *timesCustomImpl) Now(timeGMT *int) time.Time {
	location, _ := time.LoadLocation("Asia/Jakarta")

	newTimeGMT := t.gmt
	if timeGMT != nil {
		return time.Now().In(location).Add(time.Hour * time.Duration(newTimeGMT))
	}

	return time.Now().In(location)
}

func (t *timesCustomImpl) TimeStampToDateStr(timeStr, layout string) string {

	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)

	return tm.Format(layout)
}

func (t *timesCustomImpl) TimeStampToDate(timeStr, layout string) time.Time {
	i, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)

	tmStr := tm.Format(layout)

	parsed, _ := time.Parse(layout, tmStr)
	return parsed
}
