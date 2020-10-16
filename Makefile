all: gochat_client gochat_server

gochat_client: client.go
	go build -o gochat_client -v client.go

gochat_server: server.go
	go build -o gochat_server -v server.go

proto: chat.proto
	protoc --go_out=plugins=grpc:chat chat.proto 

clean:
	rm gochat_server
	rm gochat_client