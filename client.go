package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"github.com/Pilesoft/gochat/chat"
	"google.golang.org/grpc"
)

var username string
var userid int32

func main() {
	a := app.New()

	// Setup login window
	loginWindow := a.NewWindow("Chat login")
	loginWindow.Resize(fyne.NewSize(600, 200))
	loginWindow.CenterOnScreen()
	loginWindow.SetFixedSize(true)

	// Setup login form
	c := widget.NewVBox()
	addr := widget.NewEntry()
	addr.SetPlaceHolder("Server address")
	name := widget.NewEntry()
	name.SetPlaceHolder("Enter your name")
	results := widget.NewLabel("")
	btn := widget.NewButton("LOGIN", func() {
		lmsg := fmt.Sprintf("Logging into server %s with name %s\n", addr.Text, name.Text)
		fmt.Printf(lmsg)
		results.SetText(lmsg)
		conn, err := grpc.Dial(addr.Text+":9000", grpc.WithInsecure())
		if err != nil {
			msg := fmt.Sprintf("Can't connect to server: %s\n", err)
			results.SetText(msg)
			log.Print(msg)
			return
		}

		client := chat.NewChatServiceClient(conn)

		resp, err := client.Login(context.Background(), &chat.LoginRequest{Name: name.Text})
		if err != nil {
			msg := fmt.Sprintf("Login error: %s\n", err)
			results.SetText(msg)
			log.Print(msg)
			return
		}
		if resp.Status {
			msg := fmt.Sprintf("Successful login [id %d], message from server: %s\n", resp.Id, resp.Message)
			username = name.Text
			userid = resp.Id
			results.SetText(msg)
			log.Print(msg)
			time.Sleep(2 * time.Second)
			loginWindow.Close()
			startChat(&a, &client)
		} else {
			msg := fmt.Sprintf("Rejected login: %s\n", resp.Message)
			results.SetText(msg)
			log.Print(msg)
		}
	})
	c.Append(addr)
	c.Append(name)
	c.Append(btn)
	c.Append(results)

	loginWindow.SetContent(c)
	loginWindow.ShowAndRun()
}

func startChat(a *fyne.App, c *chat.ChatServiceClient) {
	title := fmt.Sprintf("Awesome gRPC chat [%s]", username)
	w := (*a).NewWindow(title)
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	names := widget.NewVBox()
	names.Append(widget.NewLabel(username))
	nameScroll := widget.NewVScrollContainer(names)

	chatWindow := widget.NewVBox()
	chatWindowScroll := widget.NewVScrollContainer(chatWindow)

	chatMessage := widget.NewMultiLineEntry()
	chatMessage.Resize(fyne.NewSize(600, 200))

	stream, err := (*c).StreamChat(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	helloMsg := chat.StreamMessage{
		Type: "HELLO",
		Name: username,
	}

	if err = stream.Send(&helloMsg); err != nil {
		log.Fatalf("Can't complete handshake with server: %v\n", err)
	}

	w.SetOnClosed(func() {
		if err := stream.CloseSend(); err != nil {
			log.Fatalln(err)
		}
	})

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("Connection closed by server")
				w.Close()
			}
			if err != nil {
				log.Fatalf("Error receiving message: %v\n", err)
			}

			switch msg.Type {
			case "CLIENTS":
				names.Children = []fyne.CanvasObject{}
				clients := strings.Split(msg.Content, ",")
				for _, client := range clients {
					names.Append(widget.NewLabel(client))
				}
				nameScroll.ScrollToBottom()

				if len(msg.Name) > 0 {
					byeStr := fmt.Sprintf("---- %s has left the chat ----", msg.Name)
					chatWindow.Append(widget.NewLabel(byeStr))
				}
			case "MESSAGE":
				messageStr := fmt.Sprintf("%s ---> %s", msg.Name, msg.Content)
				chatWindow.Append(widget.NewLabel(messageStr))
				chatWindowScroll.ScrollToBottom()
			default:
				log.Printf("Unknown message type %s: %v\n", msg.Type, msg)
			}
		}
	}()

	chatSend := widget.NewButton("Send", func() {
		fmt.Printf("We will send some cool message to the server: %s\n", chatMessage.Text)
		msg := chat.StreamMessage{
			Type:    "MESSAGE",
			Content: chatMessage.Text,
			Name:    username,
		}

		chatMessage.SetText("")

		if err := stream.Send(&msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	})

	chatBottom := widget.NewHBox(chatMessage, chatSend)

	chatArea := widget.NewVSplitContainer(chatWindowScroll, chatBottom)
	chatArea.SetOffset(10)

	genContainer := widget.NewHSplitContainer(nameScroll, chatArea)
	genContainer.SetOffset(-10)

	w.SetContent(genContainer)
	w.Show()
}
