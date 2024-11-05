package cbor2flat

import (
	"bytes"
	"context"
	"iter"
	"log"
	"slices"

	cf "github.com/takanoriyanagitani/go-cbor-flat"
)

// Flatten the nested array somehow.
type CborArrayToFlatIter func(context.Context, cf.CborArray) iter.Seq[cf.CborArray]

type CborArrToFlatIter struct {
	// Serializes the selected value(will be passed to the converter).
	AnyToSerInput

	// Selects the value to be serialized from the array.
	ArrToAny

	// Sets the converted value.
	SetAny

	// The converter(serialized CBOR -> serialized CBOR)
	SerToSer

	// Parses the converted CBOR.
	SerToArrayOut
}

func (c CborArrToFlatIter) ToSerializedToIter() SerializedToIter {
	return SerToIter{
		c.SerToSer,
		c.SerToArrayOut,
	}.ToSerializedToIter()
}

func (c CborArrToFlatIter) ToAnyToIter() AnyToIter {
	return c.ToSerializedToIter().ToAnyToIter(c.AnyToSerInput)
}

func (c CborArrToFlatIter) ToFlatIter() CborArrayToFlatIter {
	return c.ToAnyToIter().ToFlatIter(
		c.ArrToAny,
		c.SetAny,
	)
}

// Flatten the iterator of the nested array.
func (f CborArrayToFlatIter) Flatten(
	ctx context.Context,
	i cf.CborInputs,
) iter.Seq[cf.CborArray] {
	return func(yield func(cf.CborArray) bool) {
		var originals iter.Seq[cf.CborArray] = i()
		for original := range originals {
			var iflat iter.Seq[cf.CborArray] = f(ctx, original)
			for flat := range iflat {
				var flt cf.CborArray = flat
				if !yield(flt) {
					return
				}
			}
		}
	}
}

// Converts the selected value(any) to an iterator.
type AnyToIter func(context.Context, any) iter.Seq[any]

// Selects a value(any) from the array.
type ArrToAny func(cf.CborArray) any

// Creates an [ArrToAny] which gets an any using the index of the array.
func ArrToAnyByIndexNew(idx uint32) ArrToAny {
	return func(a cf.CborArray) any {
		var sz uint32 = uint32(len(a))
		var ok bool = idx < sz
		switch ok {
		case true:
			return a[idx]
		default:
			return nil
		}
	}
}

// Sets the converted value(any).
type SetAny func(cf.CborArray, any)

// Creates a [SetAny] which overwrites specified element of the array.
func SetAnyByIndexNew(idx uint32) SetAny {
	return func(a cf.CborArray, val any) {
		var sz uint32 = uint32(len(a))
		if idx < sz {
			a[idx] = val
		}
	}
}

// Creates [CborArrayToFlatIter] using [ArrToAny].
func (i AnyToIter) ToFlatIter(a2a ArrToAny, sa SetAny) CborArrayToFlatIter {
	return func(ctx context.Context, a cf.CborArray) iter.Seq[cf.CborArray] {
		var selected any = a2a(a)
		var iflat iter.Seq[any] = i(ctx, selected)
		return func(yield func(cf.CborArray) bool) {
			for ia := range iflat {
				sa(a, ia)
				if !yield(a) {
					return
				}
			}
		}
	}
}

type SerializedInput []byte

// Converts the selected value(serialized any) to an iterator.
type SerializedToIter func(context.Context, SerializedInput) iter.Seq[any]

// Serializes the selected value(any).
type AnyToSerInput func(context.Context, any) (SerializedInput, error)

// Serializes the selected value(any) into the buffer.
type AnyToSerInputBuf func(context.Context, any, *bytes.Buffer) error

func (b AnyToSerInputBuf) ToAnyToSerInput() AnyToSerInput {
	var buf bytes.Buffer
	var err error
	return func(ctx context.Context, input any) (SerializedInput, error) {
		buf.Reset()
		err = b(ctx, input, &buf)
		return buf.Bytes(), err
	}
}

func (s SerializedToIter) ToAnyToIter(a2s AnyToSerInput) AnyToIter {
	var empty []any
	return func(ctx context.Context, i any) iter.Seq[any] {
		serialized, e := a2s(ctx, i)
		if nil != e {
			return slices.Values(empty)
		}
		return s(ctx, serialized)
	}
}

type SerializedOutput []byte

// Parses the serialized array([]byte) and gets the parsed array.
type SerToArrayOut func(context.Context, SerializedOutput) (cf.CborArray, error)

// Converts the serialized any([]byte) and gets serialized array([]byte).
type SerToSer func(context.Context, SerializedInput) (SerializedOutput, error)

type SerToIter struct {
	SerToSer
	SerToArrayOut
}

func (s SerToIter) ToSerializedToIter() SerializedToIter {
	var empty []any
	return func(ctx context.Context, i SerializedInput) iter.Seq[any] {
		serout, e := s.SerToSer(ctx, i)
		if nil != e {
			log.Printf("unable to convert: %v\n", e)
			return slices.Values(empty)
		}

		arr, e := s.SerToArrayOut(ctx, serout)
		if nil != e {
			log.Printf("unable to convert to array: %v\n", e)
			return slices.Values(empty)
		}

		return slices.Values(arr)
	}
}
