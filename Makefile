build_wasm:
	GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./cmd/wasm/main.go

serv: build_wasm
	cd cmd/server; go run server.go

test_pkg:
	go test ./pkg/...

bench:
	cd ./pkg/j2s; go test -bench . -benchmem
