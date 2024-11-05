package cbor2arr

import (
	"io"
	"iter"
	"log"

	fc "github.com/fxamacker/cbor/v2"

	cf "github.com/takanoriyanagitani/go-cbor-flat"

	ic "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/cbor2arrays"
)

type CborToArrIter struct {
	*fc.Decoder
}

func (c CborToArrIter) ToIter() iter.Seq[cf.CborArray] {
	return func(yield func(cf.CborArray) bool) {
		var buf []any
		var err error
		for {
			clear(buf)
			buf = buf[:0]

			err = c.Decoder.Decode(&buf)
			if nil != err {
				if err != io.EOF {
					log.Printf("unable to decode: %v\n", err)
				}
				return
			}

			if !yield(buf) {
				return
			}
		}
	}
}

func (c CborToArrIter) AsCborArrSource() ic.CborArraySource {
	return c.ToIter
}

func CborToArrIterNew(rdr io.Reader) CborToArrIter {
	return CborToArrIter{Decoder: fc.NewDecoder(rdr)}
}
