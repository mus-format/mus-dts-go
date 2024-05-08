package dts

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go"
)

// New creates a new DTS.
func New[T any](dtm com.DTM, m mus.Marshaller[T], u mus.Unmarshaller[T],
	s mus.Sizer[T]) DTS[T] {
	return DTS[T]{dtm, m, u, s}
}

// DTS provides data type metadata (DTM) support for the mus-go serializer. It
// helps to encode DTM + data.
//
// Implements mus.Marshaller, mus.Unmarshaller and mus.Sizer interfaces.
type DTS[T any] struct {
	dtm com.DTM
	m   mus.Marshaller[T]
	u   mus.Unmarshaller[T]
	s   mus.Sizer[T]
}

// DTM returns the value with which DTS was initialized.
func (dts DTS[T]) DTM() com.DTM {
	return dts.dtm
}

// MarshalMUS marshals DTM + data.
func (dts DTS[T]) MarshalMUS(t T, bs []byte) (n int) {
	n = MarshalDTM(dts.dtm, bs)
	n += dts.m.MarshalMUS(t, bs[n:])
	return
}

// UnmarshalMUS unmarshals DTM + data.
//
// Returns ErrWrongDTM if the unmarshalled DTM differs from the dts.DTM().
func (dts DTS[T]) UnmarshalMUS(bs []byte) (t T, n int, err error) {
	dtm, n, err := UnmarshalDTM(bs)
	if err != nil {
		return
	}
	if dtm != dts.dtm {
		err = ErrWrongDTM
		return
	}
	var n1 int
	t, n1, err = dts.UnmarshalData(bs[n:])
	n += n1
	return
}

// SizeMUS calculates the size of the DTM + data.
func (dts DTS[T]) SizeMUS(t T) (size int) {
	size = SizeDTM(dts.dtm)
	return size + dts.s.SizeMUS(t)
}

// UnmarshalData unmarshals only data.
func (dts DTS[T]) UnmarshalData(bs []byte) (t T, n int, err error) {
	return dts.u.UnmarshalMUS(bs)
}
