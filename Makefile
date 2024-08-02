build:
	go build -o ./bin/smartcard-api.exe ./main.go

dev:	
	go run main.go

example:
	go run ./cmd/example/main.go


# build-wasm:
# 	GOOS=js GOARCH=wasm go build -o bin/wasm/thai-smartcard-agent.wasm ./cmd/agent/main.go
# $Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o bin/wasm/thai-smartcard-agent.wasm ./cmd/agent/main.go