package chat

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"golang.org/x/net/context"
)

// Server implementation of grcp server interface
type Server struct {
	IDCount int
	Clients map[string]*Client
}

// Client structure to hold client info
type Client struct {
	name    string
	id      int
	updated time.Time
	stream  *ChatService_StreamChatServer
}

// NewServer constructor for server type
func NewServer() *Server {
	return &Server{
		IDCount: 0,
		Clients: make(map[string]*Client),
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
	if s.Clients[msg.Name] != nil {
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

	s.Clients[msg.Name] = &c

	return &LoginResponse{
		Status:  true,
		Message: "Welcome to the chat server",
		Id:      int32(c.id),
	}, nil
}

// StreamChat process continuous chat flow
func (s *Server) StreamChat(stream ChatService_StreamChatServer) error {
	var username string

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			log.Println("Closing stream connection")
			delete(s.Clients, username)

			clientNames := []string{}
			for name := range s.Clients {
				clientNames = append(clientNames, name)
			}
			byeMsg := StreamMessage{
				Type:    "CLIENTS",
				Content: strings.Join(clientNames, ","),
				Name:    username,
			}
			for _, client := range s.Clients {
				if err := (*client.stream).Send(&byeMsg); err != nil {
					log.Printf("Error sending response to client: %v\n", err)
				}
			}
			return nil
		}
		if err != nil {
			return err
		}

		switch msg.Type {
		case "HELLO":
			username = msg.Name
			s.Clients[username].stream = &stream
			log.Printf("HELLO message received from %s\n", username)

			clientNames := []string{}
			for name := range s.Clients {
				clientNames = append(clientNames, name)
			}

			helloResponse := StreamMessage{
				Type:    "CLIENTS",
				Content: strings.Join(clientNames, ","),
			}

			for _, client := range s.Clients {
				if err := (*client.stream).Send(&helloResponse); err != nil {
					log.Printf("Error sending response to client: %v\n", err)
				}
			}
		case "MESSAGE":
			log.Printf("Received chat message from %s: %s\n", msg.Name, msg.Content)
			for _, client := range s.Clients {
				if err := (*client.stream).Send(msg); err != nil {
					log.Printf("Error broadcasting message to client: %v\n", err)
				}
			}
		default:
			log.Printf("Unknown message type %s: %v\n", msg.Type, msg)
		}
	}
}
