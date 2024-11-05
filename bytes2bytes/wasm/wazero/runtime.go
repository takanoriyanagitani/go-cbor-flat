package b2wazero2b

import (
	"context"
	"errors"

	w0 "github.com/tetratelabs/wazero"
)

var (
	ErrUnableToCompile error = errors.New("unable to compile wasm")
)

type Runtime struct {
	w0.Runtime
	w0.ModuleConfig

	InitializeInputBuffer string
	GetInputOffset        string
	Convert               string
	GetOutputOffset       string
}

func (r Runtime) Close(ctx context.Context) error {
	return r.Runtime.Close(ctx)
}

func (r Runtime) IntoCompiled(
	ctx context.Context,
	wasm []byte,
) (Compiled, error) {
	compiled, e := r.Runtime.CompileModule(ctx, wasm)
	if nil != e {
		return Compiled{}, errors.Join(ErrUnableToCompile, e, r.Close(ctx))
	}
	return Compiled{
		Runtime:        r.Runtime,
		CompiledModule: compiled,
		ModuleConfig:   r.ModuleConfig,

		InitializeInputBuffer: r.InitializeInputBuffer,
		GetInputOffset:        r.GetInputOffset,
		Convert:               r.Convert,
		GetOutputOffset:       r.GetOutputOffset,
	}, nil
}
