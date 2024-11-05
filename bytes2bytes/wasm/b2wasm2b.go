package b2wasm2b

import (
	"context"
)

type BytesToWasmToBytes func(context.Context, []byte) ([]byte, error)

type WasmSource func(context.Context) ([]byte, error)
