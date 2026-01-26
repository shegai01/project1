.phony: run
PORT=8080
run: 
	go run ./cmd/main.go

stop:
	@fuser -k ${PORT}/tcp

