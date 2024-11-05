package b2wazero2b

import (
	"context"

	wa "github.com/tetratelabs/wazero/api"

	bw "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm"
)

// Converts input bytes to bytes.
//
//   - Initialize input buffer
//   - Copy input bytes(host -> wasm)
//   - Convert
//   - Copy output bytes(wasm -> host)
type BytesToWazeroToBytes struct {
	wa.Module

	InitializeInputBuffer
	GetInputOffset
	Memory
	Convert
	GetOutputOffset
}

func (b BytesToWazeroToBytes) Close(ctx context.Context) error {
	return b.Module.Close(ctx)
}

func (b BytesToWazeroToBytes) Valid() bool {
	oks := []bool{
		nil != b.Memory.Memory,

		nil != b.InitializeInputBuffer.Function,
		nil != b.GetInputOffset.Function,
		nil != b.Convert.Function,
		nil != b.GetOutputOffset.Function,
	}

	for _, ok := range oks {
		var ng bool = !ok
		if ng {
			return false
		}
	}

	return true
}

func (b BytesToWazeroToBytes) BytesToBytes(
	ctx context.Context,
	input []byte,
) (output []byte, e error) {
	var isz uint32 = uint32(len(input))

	_, e = b.InitializeInputBuffer.InitializeDefault(ctx, isz)
	if nil != e {
		return nil, e
	}

	ioff, e := b.GetInputOffset.GetOffset(ctx)
	if nil != e {
		return nil, e
	}

	e = b.Memory.WriteBytes(ioff, input)
	if nil != e {
		return nil, e
	}

	osz, e := b.Convert.Conv(ctx)
	if nil != e {
		return nil, e
	}

	ooff, e := b.GetOutputOffset.GetOffset(ctx)
	if nil != e {
		return nil, e
	}

	return b.Memory.GetView(ooff, osz)
}

func (b BytesToWazeroToBytes) AsBytesToWasmToBytes() bw.BytesToWasmToBytes {
	return b.BytesToBytes
}
