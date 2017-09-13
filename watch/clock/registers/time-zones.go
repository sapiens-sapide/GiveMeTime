package registers

type TimeZones int8

func (tz *TimeZones) Set(value interface{}) {
	*tz = TimeZones(value.(int))
}

func (tz TimeZones) Status() interface{} {
	return int8(tz)
}
