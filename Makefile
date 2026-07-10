run:
	go run .

build:
	go build -o ./server.exe .

build-and-run:
	go build -o ./server.exe .
	server.exe