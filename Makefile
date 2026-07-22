build:
	go build -ldflags "-s -w" -o ./server.exe local/local.go

run:
	go run local/local.go

build-and-run:
	go build -ldflags "-s -w" -o ./server.exe local/local.go
	server.exe