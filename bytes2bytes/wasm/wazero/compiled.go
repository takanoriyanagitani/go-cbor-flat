package b2wazero2b

import (
	"context"
	"errors"

	w0 "github.com/tetratelabs/wazero"
)

type Compiled struct {
	w0.Runtime
	w0.CompiledModule
	w0.ModuleConfig

	InitializeInputBuffer string
	GetInputOffset        string
	Convert               string
	GetOutputOffset       string
}

func (c Compiled) Close(ctx context.Context) error {
	return errors.Join(c.CompiledModule.Close(ctx), c.Runtime.Close(ctx))
}

func (c Compiled) ToInstance(ctx context.Context) (Instance, error) {
	instance, e := c.Runtime.InstantiateModule(
		ctx,
		c.CompiledModule,
		c.ModuleConfig,
	)
	return Instance{
		Module: instance,

		InitializeInputBuffer: c.InitializeInputBuffer,
		GetInputOffset:        c.GetInputOffset,
		Convert:               c.Convert,
		GetOutputOffset:       c.GetOutputOffset,
	}, e
}
