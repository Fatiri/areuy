package manipulator

import (
	"fmt"
	"time"
)

var (
	// FormatDateTime ...
	FormatDateTimeV1 = `2006-01-02 15:04:05`
	FormatDateTimeV2 = `2006-01-02 15:04`
	// FormatDate ...
	FormatDate = `2006-01-02`
	// FormatTime ...
	FormatTime = `15:04`
	// DateMonthFormat ...
	DateMonthFormat = `02 January 2006`
	loc, _          = time.LoadLocation(`Asia/Jakarta`)

	FormatDateIdn = `02-01-2006 15:04`
)

// check date equal
func CheckDateEqual(today, newDate time.Time) bool {
	return today.Format("20060102") == newDate.Format("20060102")
}

// check time valid
func CheckTimeValid(date time.Time) error {
	today := time.Now()
	hours, minutes, _ := today.Clock()
	newHours, newMinutes, _ := date.Clock()

	todayTime := hours*60 + minutes
	dateTime := newHours*60 + newMinutes

	if todayTime > dateTime {
		return fmt.Errorf("Waktu melewati waktu perilisan")
	}

	if dateTime-todayTime < 10 {
		return fmt.Errorf("Waktu melewati batas perubahan data")
	}

	return nil
}

func CombineDateTime(date, dateTime string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", date+" "+dateTime)
}

// SetToLateNight receive date parameter value "YYYY-MM-DD"
func SetToLateNight(date string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", date)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 23, 59, 59, 0, parsed.Location()), nil
}

// CheckDateInput for validation
func CheckDateInput(date string) (bool, error) {

	tstart, err := SetToLateNight(time.Now().Format("2006-01-02"))
	if err != nil {
		return false, fmt.Errorf("cannot parse startdate: %v", err)
	}
	tend, err := SetToLateNight(date)
	if err != nil {
		return false, fmt.Errorf("cannot parse date: %v", err)
	}

	if tstart.After(tend) {
		return false, fmt.Errorf("Tanggal Expired tidak boleh kurang dari tanggal sekarang")
	}
	return true, err
}

// UTC7 ...
func UTC7(t time.Time) time.Time {
	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.Now()
	}
	return t.In(location)
}

func NowUTC7() time.Time {
	return time.Now().In(loc)
}

// ParseUTC7 ...
func ParseUTC7(timeFormat string, value string) (time.Time, error) {
	timeUTC7, err := time.ParseInLocation(timeFormat, value, loc)
	if err != nil {
		return time.Now(), err
	}

	return timeUTC7, nil
}

func StartDate(t time.Time) time.Time {
	// it will return 2009-11-10 00:00:00
	time, _ := ParseUTC7(FormatDate, t.Format(FormatDate))
	return time
}

func EndDate(t time.Time) time.Time {
	// it will return ex 2009-11-10 23:59:59
	return StartDate(t).Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
}

func StartDateString(t string) time.Time {
	// it will return 2009-11-10 00:00:00
	time, _ := ParseUTC7(FormatDate, t)
	return time
}

func EndDateString(t string) time.Time {
	// it will return ex 2009-11-10 23:59:59
	return StartDateString(t).Add(23 * time.Hour).Add(59 * time.Minute).Add(59 * time.Second)
}

func ParseTimeToMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func ParseMillisecondToTime(t int64) time.Time {
	return time.Unix(0, t*int64(time.Millisecond))
}

func ParseDateString(t string) (time.Time, error) {
	// it will return 2009-11-10 00:00
	timeParse, err := time.Parse(FormatDateIdn, t)
	if err != nil {
		return time.Now(), err
	}
	return timeParse, nil
}

func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = start.AddDate(0, 0, 6)
	return
}

func WeekStart(year, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.Local)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}
