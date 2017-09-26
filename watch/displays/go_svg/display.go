package web_svg

import (
	"bytes"
	"io"
	"sync"
)

type Display struct {
	lock sync.Mutex
	buf  bytes.Buffer
}

var Disp Display

func init() {
	Disp = Display{
		lock: sync.Mutex{},
		buf:  bytes.Buffer{},
	}
	Disp.buf.Grow(4100)
}

func (d *Display) Write(p []byte) (n int, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.buf.Write(p)
}

func (d *Display) Reset() {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.buf.Reset()
}

func (d *Display) CopyTo(w io.Writer) (n int64, err error) {
	d.lock.Lock()
	defer d.lock.Unlock()
	m, e := w.Write(d.buf.Bytes())
	return int64(m), e
}
