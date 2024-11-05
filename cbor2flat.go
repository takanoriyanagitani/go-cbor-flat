package cbor2flat

import (
	"context"
	"iter"
)

type CborArray []any

type CborOutput func(context.Context, CborArray) error

type CborInputs func() iter.Seq[CborArray]

func (o CborOutput) OutputAll(
	ctx context.Context,
	i iter.Seq[CborArray],
) error {
	for arr := range i {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		e := o(ctx, arr)
		if nil != e {
			return e
		}
	}
	return nil
}
