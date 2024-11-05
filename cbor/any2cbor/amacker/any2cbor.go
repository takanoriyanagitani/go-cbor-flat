package any2cbor

import (
	"bytes"
	"context"

	fc "github.com/fxamacker/cbor/v2"

	ca "github.com/takanoriyanagitani/go-cbor-flat/cbor/any2cbor"
)

func AnyToSerializedToBuf(_ context.Context, a any, buf *bytes.Buffer) error {
	return fc.MarshalToBuffer(a, buf)
}

var AnyToCborToBuf ca.AnyToCborToBuf = AnyToSerializedToBuf
