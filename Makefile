build_wasm:
	GOOS=js GOARCH=wasm go build -o ./web/main.wasm ./cmd/wasm/main.go

serv:
	cd cmd/server; go run server.go
