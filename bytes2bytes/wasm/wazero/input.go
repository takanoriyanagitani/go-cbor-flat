package b2wazero2b

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"

	pa "github.com/takanoriyanagitani/go-cbor-flat/util/pair"
)

var (
	ErrUnableToInitBuffer     error = errors.New("unable to init buffer")
	ErrUnableToGetInputOffset error = errors.New("unable to get input offset")
)

type InitializeInputBuffer struct{ wa.Function }
type GetInputOffset struct{ wa.Function }

func (i InitializeInputBuffer) Initialize(
	ctx context.Context,
	size uint32,
	init uint8,
) (uint32, error) {
	var esz uint64 = wa.EncodeU32(size)
	var ib8 uint64 = wa.EncodeU32(uint32(init))

	var presults pa.Pair[error, []uint64] = pa.PairNew(
		i.Function.Call(ctx, esz, ib8),
	)
	var pres pa.Pair[error, uint64] = pa.AndThen(
		presults,
		func(results []uint64) pa.Pair[error, uint64] {
			switch len(results) {
			case 1:
				return pa.Right[error, uint64](results[0])
			default:
				return pa.Left[error, uint64](ErrUnableToInitBuffer)
			}
		},
	)
	var pint32 pa.Pair[error, int32] = pa.Map(pres, wa.DecodeI32)
	var puint32 pa.Pair[error, uint32] = pa.AndThen(
		pint32,
		func(i int32) pa.Pair[error, uint32] {
			if i < 0 {
				return pa.Left[error, uint32](ErrUnableToInitBuffer)
			}
			return pa.Right[error, uint32](uint32(i))
		},
	)

	return puint32.Right, puint32.Left
}

func (i InitializeInputBuffer) InitializeDefault(
	ctx context.Context,
	size uint32,
) (uint32, error) {
	return i.Initialize(ctx, size, 0)
}

func (i GetInputOffset) GetOffset(ctx context.Context) (uint32, error) {
	var presults pa.Pair[error, []uint64] = pa.PairNew(i.Function.Call(ctx))
	var pres pa.Pair[error, uint64] = pa.AndThen(
		presults,
		func(results []uint64) pa.Pair[error, uint64] {
			switch len(results) {
			case 1:
				return pa.Right[error, uint64](results[0])
			default:
				return pa.Left[error, uint64](ErrUnableToGetInputOffset)
			}
		},
	)
	var pint32 pa.Pair[error, int32] = pa.Map(pres, wa.DecodeI32)
	var puint32 pa.Pair[error, uint32] = pa.AndThen(
		pint32,
		func(i int32) pa.Pair[error, uint32] {
			if i < 0 {
				return pa.Left[error, uint32](ErrUnableToGetInputOffset)
			}
			return pa.Right[error, uint32](uint32(i))
		},
	)

	return puint32.Right, puint32.Left
}
