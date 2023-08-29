# mus-dts-go
mus-dts-go provides DTM (**D**ata **T**ype **M**etadata) support for the 
[mus-go](https://github.com/mus-format/mus-go) serializer. With mus-dts-go we 
can encode/decode DTM + data itself.

You can find DTM in the MUS format 
[specification](https://github.com/mus-format/specification#Data-Type-Metadata).

# Tests
Test coverage is 100%.

# How To Use
```go
package main

// First of all, we have to define DTMs.
const (
  FooDTM dts.DTM = iota
  BarDTM
)

type Foo {...}

// Then Marshal/Unmarshal/Size functions.
func MarshalFooMUS(foo Foo, bs []byte) (n int) {...}

func UnmarshalFooMUS(bs []byte) (foo Foo, n int, err error) {...}

func SizeFooMUS(foo Foo) (size int) {...}

// FooDTS is created with a DTM, Marshaller, Unmarshaller and Sizer.
var FooDTS = dts.NewDTS[Foo](FooDTM, 
  mus.MarshallerFn[Foo](MarshalFooMUS),
  mus.UnmarshallerFn[Foo](UnmarshalFooMUS),
  mus.SizerFn[Foo](SizeFooMUS),
)

func main() {
  // Let's try to use FooDTS.
  var (
    foo = Foo{...}
    bs  = make([]byte, FooDTS.SizeMUS(foo))
  )
  // After marshal we are expecting to receive the following bs, where the first 
  // few bytes are FooDTM, and the rest are Foo data.
  FooDTS.MarshalMUS(foo, bs)

  // Unmarshal should return the same foo.
  afoo, _, _ := FooDTS.UnmarshalMUS(bs)
  assert.EqualDeep(afoo, foo)

  // And if the encoded DTM in bs is not FooDTM, we will receive 
  // dts.ErrWrongDTM.
  bs[0] = byte(BarDTM)
  _, _, err := FooDTS.UnmarshalMUS(bs)
  assert.EqualError(err, dts.ErrWrongDTM)
}
```

You can find the full code at [mus-examples-go](https://github.com/mus-format/mus-examples-go/tree/main/dts).