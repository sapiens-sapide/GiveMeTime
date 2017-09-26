package registers

type TimeZones float32

func (tz *TimeZones) Set(value interface{}) {
	*tz = TimeZones(value.(float32))
}

func (tz TimeZones) Status() interface{} {
	return float32(tz)
}
