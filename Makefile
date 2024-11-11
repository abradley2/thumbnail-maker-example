all: wasm-resize.wasm wasm_exec.js

wasm-resize.wasm:
	GOARCH=wasm GOOS=js go build -o wasm_resize.wasm

wasm_exec.js:
	cp "$(GOROOT)/misc/wasm/wasm_exec.js" .

clean:
	rm -f wasm_resize.wasm
	rm -f wasm_exec.js