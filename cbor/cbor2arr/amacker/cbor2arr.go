package cbor2arr

import (
	"context"

	fc "github.com/fxamacker/cbor/v2"

	cc "github.com/takanoriyanagitani/go-cbor-flat/cbor/cbor2arr"
)

func CborBytesToArrayNewBuffered() cc.CborBytesToArray {
	var buf []any
	var err error = nil
	return func(ctx context.Context, cbor []byte) ([]any, error) {
		clear(buf)
		buf = buf[:0]
		err = fc.Unmarshal(cbor, &buf)
		return buf, err
	}
}
