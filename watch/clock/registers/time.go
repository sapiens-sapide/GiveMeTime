// package date implements routines to get/set internal datetime and related registers

package registers

import (
	"time"
)

func SetTime(h, m, s int) {
	current := Utime
	Utime = time.Date(current.Year(), current.Month(), current.Day(), h, m, s, 0, time.UTC)
	Hour.Set(Utime.Hour())
	Minute.Set(Utime.Minute())
	Second.Set(Utime.Second())
}
