package arr2cbor

import (
	"context"
	"io"

	fc "github.com/fxamacker/cbor/v2"

	cf "github.com/takanoriyanagitani/go-cbor-flat"
	ia "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/arr2cbor"
)

type ArrToCborToOut struct {
	*fc.Encoder
}

func (o ArrToCborToOut) ArrToCborToOutput(
	ctx context.Context,
	arr cf.CborArray,
) error {
	return o.Encoder.Encode(arr)
}

func (o ArrToCborToOut) AsArrayToCborToOutput() ia.ArrayToCborToOutput {
	return o.ArrToCborToOutput
}

func ArrToCborToOutNew(wtr io.Writer) ArrToCborToOut {
	return ArrToCborToOut{Encoder: fc.NewEncoder(wtr)}
}
