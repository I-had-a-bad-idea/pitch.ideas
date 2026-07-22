build:
	go build -o ./server.exe local/local.go

run:
	go run local/local.go

build-and-run:
	go build -o ./server.exe local/local.go
	server.exe