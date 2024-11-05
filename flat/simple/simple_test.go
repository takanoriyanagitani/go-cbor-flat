package cbor2flat_test

import (
	"log"
	"testing"

	"context"
	"iter"
	"slices"

	cf "github.com/takanoriyanagitani/go-cbor-flat"

	ui "github.com/takanoriyanagitani/go-cbor-flat/util/iter"

	fs "github.com/takanoriyanagitani/go-cbor-flat/flat/simple"
)

func TestSimple(t *testing.T) {
	t.Parallel()

	t.Run("CborArrayToFlatIter", func(t *testing.T) {
		t.Parallel()

		t.Run("Flatten", func(t *testing.T) {
			t.Parallel()

			t.Run("empty", func(t *testing.T) {
				t.Parallel()

				var ca fs.CborArrayToFlatIter

				var emptyarr []cf.CborArray
				var ci cf.CborInputs = func() iter.Seq[cf.CborArray] {
					return slices.Values(emptyarr)
				}

				var flat iter.Seq[cf.CborArray] = ca.Flatten(
					context.Background(),
					ci,
				)

				var cnt uint64 = ui.Count(flat)
				if 0 != cnt {
					t.Fatalf("must be empty. got: %v\n", cnt)
				}
			})

			t.Run("dummy", func(t *testing.T) {
				t.Parallel()

				var arrays []cf.CborArray = []cf.CborArray{
					{"hw", 42.195, 634, true, false},
					{"hw", 3.776, 333, false, true},
				}
				var ca fs.CborArrayToFlatIter = func(
					_ context.Context,
					_ cf.CborArray,
				) iter.Seq[cf.CborArray] {
					return slices.Values(arrays)
				}

				var nested []cf.CborArray = []cf.CborArray{
					{"hh", 1, []byte{}, []any{1, 3, 5}},
					{"ih", 2, []byte{}, []any{2, 4, 6}},
				}
				var ci cf.CborInputs = func() iter.Seq[cf.CborArray] {
					return slices.Values(nested)
				}

				var flat iter.Seq[cf.CborArray] = ca.Flatten(
					context.Background(),
					ci,
				)

				var cnt uint64 = ui.Count(flat)
				if 4 != cnt {
					t.Fatalf("must be empty. got: %v\n", cnt)
				}
			})

			t.Run("CborArrToFlatIter", func(t *testing.T) {
				t.Parallel()

				ca2fi := fs.CborArrToFlatIter{
					AnyToSerInput: func(
						_ context.Context,
						_ any,
					) (fs.SerializedInput, error) {
						return nil, nil
					},

					ArrToAny: func(_ cf.CborArray) any { return 42 },
					SetAny: func(arr cf.CborArray, a any) {
						arr[3] = a
						log.Printf("after set: %v\n", arr)
					},

					SerToSer: func(
						_ context.Context,
						_ fs.SerializedInput,
					) (fs.SerializedOutput, error) {
						return nil, nil
					},

					SerToArrayOut: func(
						_ context.Context,
						_ fs.SerializedOutput,
					) (cf.CborArray, error) {
						return []any{
							1, 3, 5,
						}, nil
					},
				}

				var ca fs.CborArrayToFlatIter = ca2fi.ToFlatIter()

				var nested []cf.CborArray = []cf.CborArray{
					{"hh", 1, 42.195, []any{1, 3, 5}},
					{"ih", 2, 42.195, []any{1, 3, 5}},
				}
				var ci cf.CborInputs = func() iter.Seq[cf.CborArray] {
					return slices.Values(nested)
				}

				var flat iter.Seq[cf.CborArray] = ca.Flatten(
					context.Background(),
					ci,
				)

				next, stop := iter.Pull(flat)
				defer stop()

				expected := [6]cf.CborArray{
					[]any{"hh", 1, 42.195, 1},
					[]any{"hh", 1, 42.195, 3},
					[]any{"hh", 1, 42.195, 5},
					[]any{"ih", 2, 42.195, 1},
					[]any{"ih", 2, 42.195, 3},
					[]any{"ih", 2, 42.195, 5},
				}

				for i := 0; i < 6; i++ {
					got, ok := next()
					if !ok {
						t.Fatalf("no value got")
					}
					var exp cf.CborArray = expected[i]
					var same bool = slices.Equal(got, exp)
					if !same {
						t.Fatalf("expected %v != got %v\n", exp, got)
					}
				}
			})
		})
	})
}
