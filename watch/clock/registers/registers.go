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

// internal registers
var Utime time.Time
