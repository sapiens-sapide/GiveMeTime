// package registers declares and init global variables

package registers

import (
	"github.com/sapiens-sapide/GiveMeTime/watch/clock/relays"
	"time"
)

type Register interface {
	Set(value interface{})
	Status() interface{}
}

// registers directly used by displays
var Second *relays.Seconds
var Minute *relays.Minutes
var Hour *relays.Hours
var Weekday *relays.Weekdays
var Day *Days
var YearDay *Days
var Month *Months
var Year *Years
var YearLength *Days
var Tz *TimeZones
var Dst *DST
var SunRelay *relays.SunStep
var MoonRelay *relays.MoonStep

// internal registers
var Utime time.Time

// return a Go time from current values stored in registers.
func Now() time.Time {
	return time.Date(int(Year.Status().(uint16)), time.Month(Month.Status().(uint8)), int(Day.Status().(uint16)), int(Hour.Status().(uint8)), int(Minute.Status().(uint8)), int(Second.Status().(uint8)), 0, time.FixedZone("custom", int(Tz.Status().(int8)) * 3600))
}