package unsafepool

import (
	"errors"
	"reflect"
	"testing"
)

type basic struct {
	A bool
	B byte
	C rune
	D int
	E uint
	F int8
	G uint8
	H int16
	I uint16
	J int32
	K uint32
	L int64
	M uint64
	N float32
	O float64
	P complex64
	Q complex128
	R uintptr
	S [8]byte
	T struct {
		U string
		V *basic
	}
}

func TestResetBasic(t *testing.T) {
	var zero basic
	pool := New(reflect.TypeOf(zero))

	var b basic
	b.A = true
	b.B = 1
	b.C = 2
	b.D = 3
	b.E = 4
	b.F = 5
	b.G = 6
	b.H = 7
	b.I = 8
	b.J = 9
	b.K = 10
	b.L = 11
	b.M = 12
	b.N = 13
	b.O = 14
	b.P = 15
	b.Q = 16
	b.R = 17
	b.S[0] = 18
	b.S[1] = 19
	b.S[2] = 20
	b.S[3] = 21
	b.S[4] = 22
	b.S[5] = 23
	b.S[6] = 24
	b.S[7] = 25
	b.T.U = "26"
	b.T.V = &b

	pool.reset(&b)
	if b != zero {
		t.Fatalf("reset failed: %+v", b)
	}
}

type reference struct {
	a []byte
	b map[string]int
	c chan struct{}
	d any
	e error // like d any
	f func()
}

func TestResetReference(t *testing.T) {
	pool := New(reflect.TypeOf(reference{}))

	var r reference
	r.a = make([]byte, 8)
	r.b = make(map[string]int)
	r.c = make(chan struct{})
	r.d = "abcd"
	r.e = errors.New("test")
	r.f = func() {}

	pool.reset(&r)

	if r.a != nil {
		t.Fatalf("reset failed: r.a != nil")
	}
	if r.b != nil {
		t.Fatalf("reset failed: r.b != nil")
	}
	if r.c != nil {
		t.Fatalf("reset failed: r.c != nil")
	}
	if r.d != nil {
		t.Fatalf("reset failed: r.d != nil")
	}
	if r.e != nil {
		t.Fatalf("reset failed: r.e != nil")
	}
	if r.f != nil {
		t.Fatalf("reset failed: r.f != nil")
	}
}

type someStruct struct {
	resetted bool
}

func (ss *someStruct) Reset() {
	ss.resetted = true
}

func TestResetIface(t *testing.T) {
	pool := New(reflect.TypeOf(someStruct{}))
	var ss someStruct
	pool.reset(&ss)
	if !ss.resetted {
		t.Fatalf("reset failed: resetted is false, not call the Reset method")
	}
}

func TestGet(t *testing.T) {
	pool := New(reflect.TypeOf(basic{}))
	b := pool.Get()
	_, ok := b.(*basic)
	if !ok {
		t.Fatalf("Get returns type: %T", b)
	}
}

type notAligned struct {
	B [3]byte
}

func TestPut(t *testing.T) {
	New(reflect.TypeOf(basic{})).Put(new(basic))

	NewOf[notAligned]().Put(new(notAligned))
}
