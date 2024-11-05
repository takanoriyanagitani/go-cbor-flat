package main

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"strconv"

	w0 "github.com/tetratelabs/wazero"

	cf "github.com/takanoriyanagitani/go-cbor-flat"

	util "github.com/takanoriyanagitani/go-cbor-flat/util"
	pa "github.com/takanoriyanagitani/go-cbor-flat/util/pair"

	ca "github.com/takanoriyanagitani/go-cbor-flat/cbor/any2cbor"
	aa "github.com/takanoriyanagitani/go-cbor-flat/cbor/any2cbor/amacker"

	cc "github.com/takanoriyanagitani/go-cbor-flat/cbor/cbor2arr"
	ba "github.com/takanoriyanagitani/go-cbor-flat/cbor/cbor2arr/amacker"

	ic "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/cbor2arrays"
	ica "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/cbor2arrays/amacker"

	ia "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/arr2cbor"
	iaa "github.com/takanoriyanagitani/go-cbor-flat/cbor/iter/arr2cbor/amacker"

	bw "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm"
	ss "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm/source/stdfs"
	b0 "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm/wazero"

	ap "github.com/takanoriyanagitani/go-cbor-flat/cbor2flat/simple"
)

const (
	ArrayToAnyIndexDefault uint32 = 0
	SetAnyIndexDefault     uint32 = 0
)

type IoConfig struct {
	io.Reader
	io.Writer
}

func (i IoConfig) ToCborToArrIter() ica.CborToArrIter {
	return ica.CborToArrIterNew(i.Reader)
}

func (i IoConfig) ToCborArraySource() ic.CborArraySource {
	return i.ToCborToArrIter().AsCborArrSource()
}

func (i IoConfig) ToArrToCborToOut() iaa.ArrToCborToOut {
	return iaa.ArrToCborToOutNew(i.Writer)
}

func (i IoConfig) ToArrayToCborToOutput() ia.ArrayToCborToOutput {
	return i.ToArrToCborToOut().AsArrayToCborToOutput()
}

func (i IoConfig) ToCborInOut() ap.CborInOut {
	return ap.CborInOut{
		CborInputs: cf.CborInputs(i.ToCborArraySource()),
		CborOutput: cf.CborOutput(i.ToArrayToCborToOutput()),
	}
}

var AnyToCborToBuf ca.AnyToCborToBuf = aa.AnyToCborToBuf

var CborBytesToArr cc.CborBytesToArray = ba.CborBytesToArrayNewBuffered()

func CborSerdeConfigNew() ap.CborSerdeConfig {
	return ap.CborSerdeConfig{
		AnyToCborToBuf:   AnyToCborToBuf,
		CborBytesToArray: ba.CborBytesToArrayNewBuffered(),
	}
}

func Uint32FromStr(s string) pa.Pair[error, uint32] {
	u, e := strconv.ParseUint(s, 10, 32)
	return pa.Pair[error, uint32]{Left: e, Right: uint32(u)}
}

func StringToIntOrAlt32u(alt uint32, s string) uint32 {
	u, e := strconv.ParseUint(s, 10, 32)
	switch e {
	case nil:
		return uint32(u)
	default:
		return alt
	}
}

var StrToArrToAnyIx func(string) uint32 = util.
	Curry(StringToIntOrAlt32u)(ArrayToAnyIndexDefault)

var StrToSetAnyIx func(string) uint32 = util.
	Curry(StringToIntOrAlt32u)(SetAnyIndexDefault)

func GetEnvByKeyNew(key string) util.IO[string] {
	return func(_ context.Context) (string, error) {
		return os.Getenv(key), nil
	}
}

var WasmSourceFsStdDefaultNew util.IO[ss.FsWasmSource] = util.ComposeIO(
	GetEnvByKeyNew("ENV_WASM_MODULE_DIR_NAME"),
	ss.FsWasmSourceStdNewDefault,
)

var WasmSourceDefaultNew util.IO[bw.WasmSource] = util.ComposeIO(
	WasmSourceFsStdDefaultNew,
	func(s ss.FsWasmSource) bw.WasmSource { return s.ToWasmSource() },
)

var GetArrToAnyIxByEnvKey util.IO[uint32] = util.ComposeIO(
	GetEnvByKeyNew("ENV_ARR2ANY_IX"),
	StrToArrToAnyIx,
)

var GetStrToSetAnyIxByEnvKey util.IO[uint32] = util.ComposeIO(
	GetEnvByKeyNew("ENV_SET_ANY_IX"),
	StrToSetAnyIx,
)

var GetIxCfgFromEnv util.IO[ap.IndexConfig] = ap.IndexConfigNew(
	GetArrToAnyIxByEnvKey,
	GetStrToSetAnyIxByEnvKey,
)

type CloseAny func(context.Context) error

var CloseAnyNop CloseAny = func(_ context.Context) error { return nil }

type CloseMany []CloseAny

func (m CloseMany) ToCloseAny() CloseAny {
	return func(ctx context.Context) error {
		var err []error
		for _, cls := range m {
			err = append(err, cls(ctx))
		}
		return errors.Join(err...)
	}
}

func appNew(
	ctx context.Context,
	rdr io.Reader,
	wtr io.Writer,
	rtm w0.Runtime,
) (a ap.App, e error, cl CloseAny) {
	var closeMany CloseMany

	wsrc, e := WasmSourceDefaultNew(ctx)
	if nil != e {
		return a, e, CloseAnyNop
	}

	wasm, e := wsrc(ctx)
	if nil != e {
		return a, e, CloseAnyNop
	}

	var w0cfg b0.Config = b0.ConfigDefault
	var w0rtm b0.Runtime = w0cfg.ToRuntime(ctx, rtm)
	compiled, e := w0rtm.IntoCompiled(ctx, wasm)
	if nil != e {
		return a, e, CloseAnyNop
	}
	closeMany = append(closeMany, compiled.Close)

	instance, e := compiled.ToInstance(ctx)
	if nil != e {
		return a, e, CloseAnyNop
	}
	converter, e := instance.IntoConverter(ctx)
	if nil != e {
		return a, e, CloseAnyNop
	}
	closeMany = append(closeMany, converter.Close)

	var btwtb bw.BytesToWasmToBytes = converter.AsBytesToWasmToBytes()

	var csc ap.CborSerdeConfig = CborSerdeConfigNew()

	ic, e := GetIxCfgFromEnv(ctx)
	if nil != e {
		return a, e, CloseAnyNop
	}

	icfg := IoConfig{
		Reader: rdr,
		Writer: wtr,
	}

	var cio ap.CborInOut = icfg.ToCborInOut()

	var wcfg ap.WasmConfig = ap.WasmConfig(btwtb)

	scfg := ap.SimpleConfig{
		CborSerdeConfig: csc,
		IndexConfig:     ic,
		CborInOut:       cio,
		WasmConfig:      wcfg,
	}

	return scfg.ToApp(), nil, CloseAnyNop
}

func stdin2stdout(ctx context.Context, rtm w0.Runtime) error {
	var br io.Reader = bufio.NewReader(os.Stdin)

	var bw *bufio.Writer = bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	app, e, cls := appNew(ctx, br, bw, rtm)
	if nil != e {
		return e
	}

	defer cls(ctx)

	return app.FlattenAll(ctx)
}

func sub(ctx context.Context) error {
	return stdin2stdout(ctx, w0.NewRuntime(ctx))
}

func main() {
	e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
