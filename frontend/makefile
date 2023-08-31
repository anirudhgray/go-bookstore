.PHONY: build

run: build
	./bin/frontend
	
build:
	GOARCH=wasm GOOS=js go build -o web/app.wasm
	GOARCH=wasm GOOS=js go build -o build/web/app.wasm
	go build -o bin/frontend
