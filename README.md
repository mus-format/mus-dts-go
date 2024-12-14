# mus-dts-go

[![Go Reference](https://pkg.go.dev/badge/github.com/mus-format/mus-dts-go.svg)](https://pkg.go.dev/github.com/mus-format/mus-dts-go)
[![GoReportCard](https://goreportcard.com/badge/mus-format/mus-dts-go)](https://goreportcard.com/report/github.com/mus-format/mus-dts-go)
[![codecov](https://codecov.io/gh/mus-format/mus-dts-go/graph/badge.svg?token=VB6E8M2PFE)](https://codecov.io/gh/mus-format/mus-dts-go)

mus-dts-go provides [DTM](https://medium.com/p/21d7be309e8d) support for the 
mus-go serializer.

DTS is useful when you need to deserialize data with an unpredictable type, 
which, for example, can denote completely different types, such as `Foo` and 
`Bar`, or different versions of the same data, such as `FooV1` and `FooV2`.

DTS encode/decode DTM (which is just a number) + data itself. Thanks to DTM, one
type can be distinguished from another, let's see how:
```go
package main

import (
	"math/rand"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
)

// Define DTMs, unique DTM for each type.
const (
	FooDTM = iota + 1
	BarDTM
)

type Foo struct{...}
type Bar struct{..}

// Define Marshal/Unmarshal/Size/Skip functions.
func MarshalFooMUS(f Foo, bs []byte) (n int) {...}
func UnmarshalFooMUS(bs []byte) (f Foo, n int, err error) {...}
func SizeFooMUS(v Foo) (size int) {...}
func SkipFooMUS(bs []byte) (n int, err error) {...}

func MarshalBarMUS(b Bar, bs []byte) (n int) {...}
func UnmarshalBarMUS(bs []byte) (b Bar, n int, err error) {...}
func SizeBarMUS(b Bar) (size int) {...}
func SkipBarMUS(bs []byte) (n int, err error) {...}

// Create DTSs.
var FooDTS = dts.New[Foo](FooDTM, 
  mus.MarshallerFn[Foo](MarshalFooMUS),
  mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
  mus.SizerFn[Foo](SizeFooMUS),
  mus.Skipper(SkipFooMUS),
)
var BarDTS = dts.New[Bar](BarDTM, 
  mus.MarshallerFn[Bar](MarshalBarMUS),
  mus.UnmarshallerFn[Bar](UnmarshalBarMUS),
  mus.SizerFn[Bar](SizeBarMUS),
  mus.Skipper(SkipBarMUS),
)

func main() {
  // Make a random data ...
  bs, err := randomData()
  if err != nil {
    panic(err)
  }
  // .. and Unmarshal DTM.
  dtm, n, err := dts.UnmarshalDTM(bs)
  if err != nil {
    panic(err)
  }
  // Deserialize and process data depending on the DTM.
  switch dtm {
    case FooDTM:
      foo, _, err := FooDTS.UnmarshalData(bs[n:])
      // process foo ...
    case BarDTM:
      bar, _, err := BarDTS.UnmarshalData(bs[n:])
      // process bar ...
    default:
      panic(fmt.Sprintf("unexpected %v DTM", dtm))
  }
}

func randomData() (bs []byte) {
  // Generate random DTM
  dtm := com.DTM(rand.Intn(2) + 1)
  switch dtm {
    // Marshal Foo
    case FooDTM:
      foo := Foo{...}
      bs = make([]byte, FooDTS.Size(foo))
      FooDTS.Marshal(foo, bs)
    // Marshal Bar
    case BarDTM:
      bar := Bar{...}
      bs = make([]byte, BarDTS.Size(bar))
      BarDTS.Marshal(bar, bs)
    default:
      panic(fmt.Sprintf("unexpected %v DTM", dtm))      
  }
  return
}
```
A full example can be found at [mus-examples-go](https://github.com/mus-format/mus-examples-go/tree/main/dts)
