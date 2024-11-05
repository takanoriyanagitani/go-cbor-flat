package simple

import (
	"bytes"
	"context"
	"errors"
	"iter"

	cf "github.com/takanoriyanagitani/go-cbor-flat"

	util "github.com/takanoriyanagitani/go-cbor-flat/util"

	fl "github.com/takanoriyanagitani/go-cbor-flat/flat/simple"

	ca "github.com/takanoriyanagitani/go-cbor-flat/cbor/any2cbor"
	cc "github.com/takanoriyanagitani/go-cbor-flat/cbor/cbor2arr"

	bw "github.com/takanoriyanagitani/go-cbor-flat/bytes2bytes/wasm"
)

type IndexConfig struct {
	ArrayToAnyIndex uint32
	SetAnyIndex     uint32
}

func IndexConfigNew(
	arr2any util.IO[uint32],
	setany util.IO[uint32],
) util.IO[IndexConfig] {
	return func(ctx context.Context) (IndexConfig, error) {
		aa, ea := arr2any(ctx)
		sa, es := setany(ctx)
		return IndexConfig{
			ArrayToAnyIndex: aa,
			SetAnyIndex:     sa,
		}, errors.Join(ea, es)
	}
}

func (i IndexConfig) ToArrToAny() fl.ArrToAny {
	return fl.ArrToAnyByIndexNew(i.ArrayToAnyIndex)
}

func (i IndexConfig) ToSetAny() fl.SetAny {
	return fl.SetAnyByIndexNew(i.SetAnyIndex)
}

type CborSerdeConfig struct {
	ca.AnyToCborToBuf
	cc.CborBytesToArray
}

type CborInOut struct {
	cf.CborInputs
	cf.CborOutput
}

type WasmConfig bw.BytesToWasmToBytes

type SimpleConfig struct {
	CborSerdeConfig
	IndexConfig
	CborInOut
	WasmConfig
}

func (c SimpleConfig) ToApp() App {
	return App{
		BytesToWasmToBytes: bw.BytesToWasmToBytes(c.WasmConfig),

		AnyToCborToBuf:   c.CborSerdeConfig.AnyToCborToBuf,
		CborBytesToArray: c.CborSerdeConfig.CborBytesToArray,

		ArrToAny: c.IndexConfig.ToArrToAny(),
		SetAny:   c.IndexConfig.ToSetAny(),

		CborInputs: c.CborInOut.CborInputs,
		CborOutput: c.CborInOut.CborOutput,
	}
}

type App struct {
	bw.BytesToWasmToBytes

	ca.AnyToCborToBuf
	cc.CborBytesToArray

	fl.ArrToAny
	fl.SetAny

	cf.CborInputs
	cf.CborOutput
}

func (a App) ToAnyToSerInputBuf() fl.AnyToSerInputBuf {
	return func(ctx context.Context, i any, b *bytes.Buffer) error {
		return a.AnyToCborToBuf(ctx, i, b)
	}
}

func (a App) ToSerToArrayOut() fl.SerToArrayOut {
	return func(
		ctx context.Context,
		o fl.SerializedOutput,
	) (cf.CborArray, error) {
		return a.CborBytesToArray(ctx, o)
	}
}

func (a App) ToCborArrToFlatIter() fl.CborArrToFlatIter {
	return fl.CborArrToFlatIter{
		AnyToSerInput: a.ToAnyToSerInputBuf().ToAnyToSerInput(),
		ArrToAny:      a.ArrToAny,
		SetAny:        a.SetAny,
		SerToSer: func(
			ctx context.Context,
			i fl.SerializedInput,
		) (fl.SerializedOutput, error) {
			return a.BytesToWasmToBytes(ctx, i)
		},
		SerToArrayOut: a.ToSerToArrayOut(),
	}
}

func (a App) ToCborArrayToFlatIter() fl.CborArrayToFlatIter {
	return a.ToCborArrToFlatIter().ToFlatIter()
}

func (a App) Flatten(ctx context.Context) iter.Seq[cf.CborArray] {
	return a.
		ToCborArrayToFlatIter().
		Flatten(ctx, a.CborInputs)
}

func (a App) FlattenAll(ctx context.Context) error {
	return a.CborOutput.OutputAll(
		ctx,
		a.Flatten(ctx),
	)
}
