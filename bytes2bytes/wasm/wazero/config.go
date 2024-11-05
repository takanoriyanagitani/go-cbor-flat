package b2wazero2b

import (
	"context"

	w0 "github.com/tetratelabs/wazero"
)

type Config struct {
	w0.ModuleConfig

	InitializeInputBuffer string
	GetInputOffset        string
	Convert               string
	GetOutputOffset       string
}

func (c Config) ToRuntime(
	ctx context.Context,
	rtm w0.Runtime,
) Runtime {
	return Runtime{
		Runtime:      rtm,
		ModuleConfig: c.ModuleConfig,

		InitializeInputBuffer: c.InitializeInputBuffer,
		GetInputOffset:        c.GetInputOffset,
		Convert:               c.Convert,
		GetOutputOffset:       c.GetOutputOffset,
	}
}

var ConfigDefault Config = Config{
	ModuleConfig: w0.NewModuleConfig().WithName(""),

	InitializeInputBuffer: "initialize_input_buffer",
	GetInputOffset:        "get_input_offset",
	Convert:               "convert",
	GetOutputOffset:       "get_output_offset",
}
