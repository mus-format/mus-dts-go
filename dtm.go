package dts

import (
	com "github.com/mus-format/common-go"
	"github.com/mus-format/mus-go/varint"
)

func MarshalDTM(dtm com.DTM, bs []byte) (n int) {
	return varint.MarshalInt(int(dtm), bs)
}

func UnmarshalDTM(bs []byte) (dtm com.DTM, n int, err error) {
	num, n, err := varint.UnmarshalInt(bs)
	if err != nil {
		return
	}
	dtm = com.DTM(num)
	return
}

func SizeDTM(dtm com.DTM) (size int) {
	return varint.SizeInt(int(dtm))
}
