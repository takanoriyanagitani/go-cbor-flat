package b2wazero2b

import (
	"context"
	"errors"

	wa "github.com/tetratelabs/wazero/api"
)

type Instance struct {
	wa.Module

	InitializeInputBuffer string
	GetInputOffset        string
	Convert               string
	GetOutputOffset       string
}

var (
	ErrInvalidInstance error = errors.New("invalid instance")
)

func (i Instance) Close(ctx context.Context) error {
	return i.Module.Close(ctx)
}

func (i Instance) IntoConverter(
	ctx context.Context,
) (BytesToWazeroToBytes, error) {
	b2w2b := BytesToWazeroToBytes{
		Module: i.Module,

		Memory: Memory{Memory: i.Module.Memory()},

		InitializeInputBuffer: InitializeInputBuffer{
			Function: i.Module.ExportedFunction(i.InitializeInputBuffer),
		},
		GetInputOffset: GetInputOffset{
			Function: i.Module.ExportedFunction(i.GetInputOffset),
		},
		Convert: Convert{
			Function: i.Module.ExportedFunction(i.Convert),
		},
		GetOutputOffset: GetOutputOffset{
			Function: i.Module.ExportedFunction(i.GetOutputOffset),
		},
	}

	var valid bool = b2w2b.Valid()
	var ng bool = !valid
	if ng {
		return BytesToWazeroToBytes{}, errors.Join(
			ErrInvalidInstance,
			i.Close(ctx),
		)
	}

	return b2w2b, nil
}
