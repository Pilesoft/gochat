all: client server

client: client.go
	go build -v client.go

server: server.go
	go build -v server.go

proto: chat.proto
	protoc --go_out=plugins=grpc:chat chat.proto 

clean:
	rm server
	rm client