package b2wazero2b

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"

	pa "github.com/takanoriyanagitani/go-cbor-flat/util/pair"
)

var (
	ErrUnableToGetOutputOffset error = errors.New("unable to get output offset")
)

type GetOutputOffset struct{ wa.Function }

func (i GetOutputOffset) GetOffset(ctx context.Context) (uint32, error) {
	var presults pa.Pair[error, []uint64] = pa.PairNew(i.Function.Call(ctx))
	var pres pa.Pair[error, uint64] = pa.AndThen(
		presults,
		func(results []uint64) pa.Pair[error, uint64] {
			switch len(results) {
			case 1:
				return pa.Right[error, uint64](results[0])
			default:
				return pa.Left[error, uint64](ErrUnableToGetOutputOffset)
			}
		},
	)
	var pint32 pa.Pair[error, int32] = pa.Map(pres, wa.DecodeI32)
	var puint32 pa.Pair[error, uint32] = pa.AndThen(
		pint32,
		func(i int32) pa.Pair[error, uint32] {
			if i < 0 {
				return pa.Left[error, uint32](ErrUnableToGetOutputOffset)
			}
			return pa.Right[error, uint32](uint32(i))
		},
	)

	return puint32.Right, puint32.Left
}
