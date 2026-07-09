all:
	go build -o ./server.exe ./cmd/server
	server.exe