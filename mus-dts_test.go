package dts

import (
	"reflect"
	"testing"

	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/varint"
)

const FooDTM com.DTM = 1

type Foo struct {
	num int
	str string
}

func MarshalFooMUS(foo Foo, bs []byte) (n int) {
	n = varint.MarshalInt(foo.num, bs)
	n += ord.MarshalString(foo.str, nil, bs[n:])
	return
}

func UnmarshalFooMUS(bs []byte) (foo Foo, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var n1 int
	foo.str, n1, err = ord.UnmarshalString(nil, bs[n:])
	n += n1
	return
}

func SizeFooMUS(foo Foo) (size int) {
	size = varint.SizeInt(foo.num)
	return size + ord.SizeString(foo.str, nil)
}

func TestDTS(t *testing.T) {

	t.Run("MarshalMUS, UnmarshalMUS, SizeMUS methods should work correctly",
		func(t *testing.T) {
			var (
				foo    = Foo{num: 11, str: "hello world"}
				fooDTS = New[Foo](FooDTM,
					mus.MarshallerFn[Foo](MarshalFooMUS),
					mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
					mus.SizerFn[Foo](SizeFooMUS),
				)
				bs = make([]byte, fooDTS.SizeMUS(foo))
			)
			n := fooDTS.MarshalMUS(foo, bs)
			if n != len(bs) {
				t.Fatalf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			afoo, n, err := fooDTS.UnmarshalMUS(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != len(bs) {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
		})

	t.Run("MarshalMUS, UnmarshalDTM, UnmarshalData, SizeMUS methods should work correctly",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				foo        = Foo{num: 11, str: "hello world"}
				fooDTS     = New[Foo](FooDTM,
					mus.MarshallerFn[Foo](MarshalFooMUS),
					mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
					mus.SizerFn[Foo](SizeFooMUS),
				)
				bs = make([]byte, fooDTS.SizeMUS(foo))
			)
			n := fooDTS.MarshalMUS(foo, bs)
			if n != len(bs) {
				t.Fatalf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			dtm, n, err := UnmarshalDTM(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", 1, n)
			}
			if dtm != FooDTM {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", FooDTM, dtm)
			}
			afoo, n, err := fooDTS.UnmarshalData(bs[n:])
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != len(bs)-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			fooDTS = New[Foo](FooDTM, nil, nil, nil)
			dtm    = fooDTS.DTM()
		)
		if dtm != FooDTM {
			t.Errorf("unexpected dtm, want '%v' actual '%v'", FooDTM, dtm)
		}
	})

	t.Run("UnamrshalMUS should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize  = 1
				bs          = []byte{byte(FooDTM) + 3}
				fooDTS      = New[Foo](FooDTM, nil, nil, nil)
				foo, n, err = fooDTS.UnmarshalMUS(bs)
			)
			if err != ErrWrongDTM {
				t.Errorf("unexpected error, want '%v' actual '%v'", ErrWrongDTM, err)
			}
			if !reflect.DeepEqual(foo, Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", wantDTSize, n)
			}
		})

	t.Run("If UnmarshalDTM fails with an error, UnmarshalMUS should return it",
		func(t *testing.T) {
			var (
				bs          = []byte{}
				fooDTS      = New[Foo](FooDTM, nil, nil, nil)
				foo, n, err = fooDTS.UnmarshalMUS(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if !reflect.DeepEqual(foo, Foo{}) {
				t.Errorf("unexpected foo, want '%v' actual '%v'", nil, foo)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If varint.UnmarshalInt fails with an error, UnmarshalDTM should return it",
		func(t *testing.T) {
			var (
				bs          = []byte{}
				dtm, n, err = UnmarshalDTM(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if dtm != 0 {
				t.Errorf("unexpected dtm, want '%v' actual '%v'", 0, dtm)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

}
