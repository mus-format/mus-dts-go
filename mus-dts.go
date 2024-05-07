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

// DTS provides data type metadata support for the mus-go serializer.
//
// It implements mus.Marshaller, mus.Unmarshaller and mus.Sizer interfaces.
type DTS[T any] struct {
	dtm com.DTM

	m mus.Marshaller[T]
	u mus.Unmarshaller[T]
	s mus.Sizer[T]
}

// DTM returns a data type metadata.
func (dts DTS[T]) DTM() com.DTM {
	return dts.dtm
}

// MarshalMUS marshals DTM and data to the MUS format.
func (dts DTS[T]) MarshalMUS(t T, bs []byte) (n int) {
	n = MarshalDTM(dts.dtm, bs)
	n += dts.m.MarshalMUS(t, bs[n:])
	return
}

// UnmarshalMUS unmarshals DTM and data from the MUS format.
//
// Returns ErrWrongDTM if DTM from bs is different.
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

// SizeMUS calculates the DTM and data size in the MUS format.
func (dts DTS[T]) SizeMUS(t T) (size int) {
	size = SizeDTM(dts.dtm)
	return size + dts.s.SizeMUS(t)
}

// UnmarshalMUS unmarshals data without DTM from the MUS format.
func (dts DTS[T]) UnmarshalData(bs []byte) (t T, n int, err error) {
	return dts.u.UnmarshalMUS(bs)
}
