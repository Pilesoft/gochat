all: gochat_client gochat_server

gochat_client: client.go
	go build -o gochat_client -v client.go

gochat_server: server.go
	go build -o gochat_server -v server.go

proto: chat.proto
	protoc --go_out=plugins=grpc:chat chat.proto

win:
	CC=/usr/bin/x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o gochat_client.exe -v client.go
	CC=/usr/bin/x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -o gochat_server.exe -v server.go

clean:
	rm gochat_server
	rm gochat_client