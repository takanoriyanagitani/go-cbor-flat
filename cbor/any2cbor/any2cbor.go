package any2cbor

import (
	"bytes"
	"context"
)

// Serializes the value(any) into the buffer.
type AnyToCborToBuf func(context.Context, any, *bytes.Buffer) error
