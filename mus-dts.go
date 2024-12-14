package dts

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
)

// New creates a new DTS.
func New[T any](dtm com.DTM, m mus.Marshaller[T], u mus.Unmarshaller[T],
	s mus.Sizer[T],
	sk mus.Skipper,
) DTS[T] {
	return DTS[T]{dtm, m, u, s, sk}
}

// DTS provides data type metadata (DTM) support for the mus-go serializer. It
// helps to encode/decode DTM + data.
//
// Implements mus.Marshaller, mus.Unmarshaller, mus.Sizer and mus.Skipper
// interfaces.
type DTS[T any] struct {
	dtm com.DTM
	m   mus.Marshaller[T]
	u   mus.Unmarshaller[T]
	s   mus.Sizer[T]
	sk  mus.Skipper
}

// DTM returns the value used to initialize DTS.
func (d DTS[T]) DTM() com.DTM {
	return d.dtm
}

// Marshal marshals DTM + data.
func (d DTS[T]) Marshal(t T, bs []byte) (n int) {
	n = MarshalDTM(d.dtm, bs)
	n += d.m.Marshal(t, bs[n:])
	return
}

// Unmarshal unmarshals DTM + data.
//
// Returns ErrWrongDTM if the unmarshalled DTM differs from the d.DTM().
func (d DTS[T]) Unmarshal(bs []byte) (t T, n int, err error) {
	dtm, n, err := UnmarshalDTM(bs)
	if err != nil {
		return
	}
	if dtm != d.dtm {
		err = ErrWrongDTM
		return
	}
	var n1 int
	t, n1, err = d.UnmarshalData(bs[n:])
	n += n1
	return
}

// Size calculates the size of the DTM + data.
func (d DTS[T]) Size(t T) (size int) {
	size = SizeDTM(d.dtm)
	return size + d.s.Size(t)
}

// Skip skips DTM + data.
//
// Returns ErrWrongDTM if the unmarshalled DTM differs from the d.DTM().
func (d DTS[T]) Skip(bs []byte) (n int, err error) {
	dtm, n, err := UnmarshalDTM(bs)
	if err != nil {
		return
	}
	if dtm != d.dtm {
		err = ErrWrongDTM
		return
	}
	var n1 int
	n1, err = d.sk.Skip(bs[n:])
	n += n1
	return
}

// UnmarshalData unmarshals only data.
func (d DTS[T]) UnmarshalData(bs []byte) (t T, n int, err error) {
	return d.u.Unmarshal(bs)
}

// SkipData skips only data.
func (d DTS[T]) SkipData(bs []byte) (n int, err error) {
	return d.sk.Skip(bs)
}
