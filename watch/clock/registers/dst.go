package registers

type DST bool

func (dst *DST) Set(value interface{}) {
	*dst = DST(value.(bool))
}

func (dst DST) Status() interface{} {
	return dst
}
