# mus-dts-go

[![Go Reference](https://pkg.go.dev/badge/github.com/mus-format/mus-dts-go.svg)](https://pkg.go.dev/github.com/mus-format/mus-dts-go)
[![GoReportCard](https://goreportcard.com/badge/mus-format/mus-dts-go)](https://goreportcard.com/report/github.com/mus-format/mus-dts-go)
[![codecov](https://codecov.io/gh/mus-format/mus-dts-go/graph/badge.svg?token=VB6E8M2PFE)](https://codecov.io/gh/mus-format/mus-dts-go)

mus-dts-go provides DTM support for the mus-go serializer. It allows to create
DTS (DTM Support) for a type.

DTSs are useful when there is a need to deserialize data, but we don't know in 
advance what type it is. For example, it could be `Foo` or `Bar` (or it could be
just different versions of the same data, like `FooV1` or `FooV2`), we just 
don't know, but want to handle both of these cases.

DTS encode/decode DTM (which is just a number) + data itself. Thanks to DTM, we 
can distinguish one type of data from another, let's see how:
```go
package main

import (
	"math/rand"

	com "github.com/mus-format/common-go"
	dts "github.com/mus-format/mus-dts-go"
	"github.com/mus-format/mus-go"
)

// First of all, we have to define DTMs, unique DTM for each type.
const (
	FooDTM = iota + 1
	BarDTM
)

type Foo struct{...}
type Bar struct{..}

// Then define Marshal/Unmarshal/Size/Skip functions.
func MarshalFooMUS(f Foo, bs []byte) (n int) {...}
func UnmarshalFooMUS(bs []byte) (f Foo, n int, err error) {...}
func SizeFooMUS(v Foo) (size int) {...}
func SkipFooMUS(bs []byte) (n int, err error) {...}

func MarshalBarMUS(b Bar, bs []byte) (n int) {...}
func UnmarshalBarMUS(bs []byte) (b Bar, n int, err error) {...}
func SizeBarMUS(b Bar) (size int) {...}
func SkipBarMUS(bs []byte) (n int, err error) {...}

// And only now we can define DTSs.
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
  // Let's make a random data
  bs, err := randomData()
  if err != nil {
    panic(err)
  }
  // and Unmarshal DTM.
  dtm, n, err := dts.UnmarshalDTM(bs)
  if err != nil {
    panic(err)
  }
  // Now we can deserialize and process data depending on the DTM.
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

# Tests
Test coverage is 100%.
