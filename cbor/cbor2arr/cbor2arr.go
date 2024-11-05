package cbor2arr

import (
	"context"
)

type CborBytesToArray func(context.Context, []byte) ([]any, error)
