package fs2wasm

import (
	"bufio"
	"context"
	"io"
	"io/fs"
	"os"

	bw "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm"
)

type FsWasmSource struct {
	fs.FS
	WasmModuleBasename string
	WasmMaxBytes       uint32
}

func (f FsWasmSource) ToWasmSource() bw.WasmSource {
	return func(_ context.Context) ([]byte, error) {
		file, e := f.FS.Open(f.WasmModuleBasename)
		if nil != e {
			return nil, e
		}
		defer file.Close()

		var br io.Reader = bufio.NewReader(file)
		limited := &io.LimitedReader{
			R: br,
			N: int64(f.WasmMaxBytes),
		}
		return io.ReadAll(limited)
	}
}

type FsConfig struct {
	WasmModuleBasename string
	WasmMaxBytes       uint32
}

var FsConfigDefault FsConfig = FsConfig{
	WasmModuleBasename: "flatten.wasm",
	WasmMaxBytes:       16777216, // 16 MiB
}

func FsWasmSourceStdNew(dirname string, cfg FsConfig) FsWasmSource {
	return FsWasmSource{
		FS:                 os.DirFS(dirname),
		WasmModuleBasename: cfg.WasmModuleBasename,
		WasmMaxBytes:       cfg.WasmMaxBytes,
	}
}

func FsWasmSourceStdNewDefault(dirname string) FsWasmSource {
	return FsWasmSourceStdNew(dirname, FsConfigDefault)
}
