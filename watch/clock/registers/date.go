package registers

import (
	"time"
)

func SetDate(y, m, d int) {
	// rely on time.Date to handle date normalization
	Utime = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC)
	Year.Set(Utime.Year())
	Month.Set(int(Utime.Month()))
	Day.Set(Utime.Day())
	YearDay.Set(Utime.YearDay())
	if Year.IsLeapYear() {
		YearLength.Set(366)
	} else {
		YearLength.Set(365)
	}
	wd := int(Utime.Weekday())
	if wd == 0 {
		Weekday.Set(7)
	} else {
		Weekday.Set(wd)
	}
}
