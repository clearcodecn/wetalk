package bool

import (
	"sync/atomic"
)

type Bool struct {
	v uint32
}

const (
	False uint32 = iota
	True
)

func (b *Bool) True() bool {
	return atomic.LoadUint32(&b.v) == True
}

func (b *Bool) SetValue(v uint32) {
	atomic.StoreUint32(&b.v, v)
}
