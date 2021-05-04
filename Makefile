all: php_form_wasm wasm_exec_js

php_form_wasm:
	GOOS=js GOARCH=wasm go build -o ./build/php-form.wasm

wasm_exec_js:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" ./build