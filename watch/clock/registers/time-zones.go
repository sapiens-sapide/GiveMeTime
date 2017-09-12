package registers

type TimeZones uint8

func (tz *TimeZones) Set(value interface{}) {
	*tz = TimeZones(value.(int))
}

func (tz TimeZones) Status() interface{} {
	return tz
}
