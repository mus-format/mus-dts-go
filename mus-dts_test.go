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

func MarshalFoo(foo Foo, bs []byte) (n int) {
	n = varint.MarshalInt(foo.num, bs)
	n += ord.MarshalString(foo.str, nil, bs[n:])
	return
}

func UnmarshalFoo(bs []byte) (foo Foo, n int, err error) {
	foo.num, n, err = varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	var n1 int
	foo.str, n1, err = ord.UnmarshalString(nil, bs[n:])
	n += n1
	return
}

func SizeFoo(foo Foo) (size int) {
	size = varint.SizeInt(foo.num)
	return size + ord.SizeString(foo.str, nil)
}

func SkipFoo(bs []byte) (n int, err error) {
	n, err = varint.SkipInt(bs)
	if err != nil {
		return
	}
	var n1 int
	n1, err = ord.SkipString(nil, bs[n:])
	n += n1
	return
}

func TestDTS(t *testing.T) {

	t.Run("Marshal, Unmarshal, Size, Skip methods should work correctly",
		func(t *testing.T) {
			var (
				foo    = Foo{num: 11, str: "hello world"}
				fooDTS = New[Foo](FooDTM,
					mus.MarshallerFn[Foo](MarshalFoo),
					mus.UnmarshallerFn[Foo](UnmarshalFoo),
					mus.SizerFn[Foo](SizeFoo),
					mus.SkipperFn(SkipFoo),
				)
				bs = make([]byte, fooDTS.Size(foo))
			)
			n := fooDTS.Marshal(foo, bs)
			if n != len(bs) {
				t.Fatalf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			afoo, n, err := fooDTS.Unmarshal(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n != len(bs) {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			n1, err := fooDTS.Skip(bs)
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != n {
				t.Errorf("unexpected n1, want '%v' actual '%v'", n, n1)
			}

		})

	t.Run("Marshal, UnmarshalDTM, UnmarshalData, Size, SkipData methods should work correctly",
		func(t *testing.T) {
			var (
				wantDTSize = 1
				foo        = Foo{num: 11, str: "hello world"}
				fooDTS     = New[Foo](FooDTM,
					mus.MarshallerFn[Foo](MarshalFoo),
					mus.UnmarshallerFn[Foo](UnmarshalFoo),
					mus.SizerFn[Foo](SizeFoo),
					mus.SkipperFn(SkipFoo),
				)
				bs = make([]byte, fooDTS.Size(foo))
			)
			n := fooDTS.Marshal(foo, bs)
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
			afoo, n1, err := fooDTS.UnmarshalData(bs[n:])
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != len(bs)-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n1)
			}
			if !reflect.DeepEqual(afoo, foo) {
				t.Errorf("unexpected afoo, want '%v' actual '%v'", foo, afoo)
			}
			n1, err = fooDTS.SkipData(bs[n:])
			if err != nil {
				t.Errorf("unexpected error, want '%v' actual '%v'", nil, err)
			}
			if n1 != len(bs)-wantDTSize {
				t.Errorf("unexpected n, want '%v' actual '%v'", len(bs), n1)
			}
		})

	t.Run("DTM method should return correct DTM", func(t *testing.T) {
		var (
			fooDTS = New[Foo](FooDTM, nil, nil, nil, nil)
			dtm    = fooDTS.DTM()
		)
		if dtm != FooDTM {
			t.Errorf("unexpected dtm, want '%v' actual '%v'", FooDTM, dtm)
		}
	})

	t.Run("Unamrshal should fail with ErrWrongDTM, if meets another DTM",
		func(t *testing.T) {
			var (
				wantDTSize  = 1
				bs          = []byte{byte(FooDTM) + 3}
				fooDTS      = New[Foo](FooDTM, nil, nil, nil, nil)
				foo, n, err = fooDTS.Unmarshal(bs)
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

	t.Run("If UnmarshalDTM fails with an error, Unmarshal should return it",
		func(t *testing.T) {
			var (
				bs          = []byte{}
				fooDTS      = New[Foo](FooDTM, nil, nil, nil, nil)
				foo, n, err = fooDTS.Unmarshal(bs)
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

	t.Run("If SkipDTM fails with an error, Skip should return it",
		func(t *testing.T) {
			var (
				bs     = []byte{}
				fooDTS = New[Foo](FooDTM, nil, nil, nil, nil)
				n, err = fooDTS.Skip(bs)
			)
			if err != mus.ErrTooSmallByteSlice {
				t.Errorf("unexpected error, want '%v' actual '%v'",
					mus.ErrTooSmallByteSlice,
					err)
			}
			if n != 0 {
				t.Errorf("unexpected n, want '%v' actual '%v'", 0, n)
			}
		})

	t.Run("If varint.UnmarshalPositiveInt fails with an error, UnmarshalDTM should return it",
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
