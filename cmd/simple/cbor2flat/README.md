## Example

input(jsonl)

```
["HW",3.776,42,["HH",273.15,333,false,true],true,false]
["hw",42.195,634,["hh",2.99792458,null,false,true],true,false]
["hw",42.195,null,["hh",2.99792458,null,false,true],true,false]
```

commands(*)

```
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
```

output(converted to jsonl)

```
["HW", 3.776, 42.0, "HH", true, false]
["HW", 3.776, 42.0, 273.15, true, false]
["HW", 3.776, 42.0, 333.0, true, false]
["HW", 3.776, 42.0, false, true, false]
["HW", 3.776, 42.0, true, true, false]
["hw", 42.195, 634.0, "hh", true, false]
["hw", 42.195, 634.0, 2.99792458, true, false]
["hw", 42.195, 634.0, null, true, false]
["hw", 42.195, 634.0, false, true, false]
["hw", 42.195, 634.0, true, true, false]
["hw", 42.195, null, "hh", true, false]
["hw", 42.195, null, 2.99792458, true, false]
["hw", 42.195, null, null, true, false]
["hw", 42.195, null, false, true, false]
["hw", 42.195, null, true, true, false]
```


(*) Requirements
- json2arr2cbor
  github.com/takanoriyanagitani/go-json2cbor/tree/main/cmd/json2arr2cbor
- python
- uv(python lib)
- cbor2(python lib)
