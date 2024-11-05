#!/bin/sh

export ENV_WASM_MODULE_DIR_NAME=./modules.d/out.d
export ENV_ARR2ANY_IX=3
export ENV_SET_ANY_IX=${ENV_ARR2ANY_IX}

cat ./sample.d/input.jsonl |
	json2arr2cbor |
	./cbor2flat |
	python3 \
		-m uv \
		tool \
		run \
		cbor2 \
		--sequence 
