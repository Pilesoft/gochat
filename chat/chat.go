package chat

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/context"
)

// Server implementation of grcp server interface
type Server struct {
	IDCount int
	Clients map[string]Client
}

// Client structure to hold client info
type Client struct {
	name    string
	id      int
	updated time.Time
}

// NewServer constructor for server type
func NewServer() *Server {
	return &Server{
		IDCount: 0,
		Clients: make(map[string]Client),
	}
}

// SayHello implementation of greeting function
func (s *Server) SayHello(ctx context.Context, msg *Message) (*Message, error) {
	log.Printf("Received message from client: %s\n", msg.Body)
	return &Message{Body: "Hello from the server!!"}, nil
}

// Login implementation of client login
func (s *Server) Login(ctx context.Context, msg *LoginRequest) (*LoginResponse, error) {
	log.Printf("Received login from user { %s }", msg.Name)
	if s.Clients[msg.Name].name == msg.Name {
		return &LoginResponse{
			Status:  false,
			Message: fmt.Sprintf("Name %s already in use", msg.Name),
		}, nil
	}

	s.IDCount++
	c := Client{
		name:    msg.Name,
		id:      s.IDCount,
		updated: time.Now(),
	}

	s.Clients[msg.Name] = c

	return &LoginResponse{
		Status:  true,
		Message: "Welcome to the chat server",
		Id:      int32(c.id),
	}, nil
}
