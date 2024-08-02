SMC_AGENT_PORT=9898
SMC_SHOW_IMAGE=true
SMC_SHOW_NHSO=true
SMC_SHOW_LASER=true

export SMC_PORT
export SMC_SHOW_IMAGE
export SMC_SHOW_NHSO
export SMC_SHOW_LASER

build:
	go build -o ./bin/smartcard-api.exe ./main.go

dev:	
	go run main.go

example:
	go run ./cmd/example/main.go

build-linux:
	go build -o ./bin/thai-smartcard-agent.linux-amd64 ./cmd/agent/main.go
	tar -czvf ./bin/thai-smartcard-agent.linux-amd64.tar.gz ./bin/thai-smartcard-agent.linux-amd64

build-mac:
	go build -o ./bin/thai-smartcard-agent.darwin-amd64 ./cmd/agent/main.go

build-win:	
#	go build -o ./bin/thai-smartcard-agent.windows-amd64.exe ./cmd/agent/main.go
# go build -o ./main.exe ./cmd/agent/main.go

	$Env:SMC_SHOW_IMAGE='true'; 
	$Env:SMC_SHOW_NHSO='true'; 
	$Env:SMC_SHOW_LASER='true'; 
	$Env:SMC_AGENT_PORT=9898; 
	$Env:API_KEY='yNJvodfppC:jgUbizp%LzfJI6YLcPbQ'; 
	go build -o ./main.exe ./cmd/agent/main.go

build-wasm:
	GOOS=js GOARCH=wasm go build -o bin/wasm/thai-smartcard-agent.wasm ./cmd/agent/main.go

# build-wasm:
# 	GOOS=js GOARCH=wasm go build -o bin/wasm/thai-smartcard-agent.wasm ./cmd/agent/main.go
# $Env:GOOS = "js"; $Env:GOARCH = "wasm"; go build -o bin/wasm/thai-smartcard-agent.wasm ./cmd/agent/main.go