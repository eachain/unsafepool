package unsafepool

import (
	"reflect"
	"sync"
	"unsafe"
)

// Pool 对sync.Pool+reflect+unsafe的封装。
// 一般用于需要使用reflect，又需要sync.Pool的情况下。
type Pool struct {
	size  int
	typ   reflect.Type
	reset func(any)
	pool  sync.Pool
}

func (p *Pool) new() any {
	return reflect.New(p.typ).Interface()
}

// Get 返回一个指针。
func (p *Pool) Get() any {
	return p.pool.Get()
}

type resetter interface {
	Reset()
}

// Put 放回一个指针，并将指针指向内存清0。
func (p *Pool) Put(x any) {
	p.reset(x)
	p.pool.Put(x)
}

func (p *Pool) selfReset(x any) {
	x.(resetter).Reset()
}

// unsafeReset 将指针指向内存清0。
func (p *Pool) unsafeReset(x any) {
	pointer := reflect.ValueOf(x).UnsafePointer()
	u64 := p.size / 8
	if u64 > 0 {
		s := unsafe.Slice((*uint64)(pointer), u64)
		for i := range s {
			s[i] = 0
		}
	}

	reset := u64 * 8
	if reset < p.size {
		s := unsafe.Slice((*uint8)(unsafe.Add(pointer, reset)), p.size-reset)
		for i := range s {
			s[i] = 0
		}
	}
}

var resetterType = reflect.TypeOf(new(resetter)).Elem()

// New返回一个子类型为*typ的内存池。
// 如果typ为int，pool.Get()返回*int。
func New(typ reflect.Type) *Pool {
	p := &Pool{
		size: int(typ.Size()),
		typ:  typ,
	}
	p.pool.New = p.new

	if reflect.PointerTo(typ).Implements(resetterType) {
		p.reset = p.selfReset
	} else {
		p.reset = p.unsafeReset
	}

	return p
}

// NewOf返回一个子类型为*T的内存池。
func NewOf[T any]() *Pool {
	return New(reflect.TypeOf(new(T)).Elem())
}
