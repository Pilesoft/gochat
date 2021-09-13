package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	//"io"
	"log"
	//"strings"
	"time"

	/*"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"*/

	g "github.com/AllenDang/giu"
	"github.com/Pilesoft/gochat/chat"
	"google.golang.org/grpc"
)

var username string
var serverAddress string
var userid int32
var logged bool

var wnd *g.MasterWindow

var client chat.ChatServiceClient
var clients []interface{}

func login() {
	lmsg := fmt.Sprintf("Logging into server %s with name %s\n", serverAddress, username)
	fmt.Printf(lmsg)
	conn, err := grpc.Dial(serverAddress+":9000", grpc.WithInsecure())
	if err != nil {
		msg := fmt.Sprintf("Can't connect to server: %s\n", err)
		log.Print(msg)
		return
	}

	client = chat.NewChatServiceClient(conn)

	resp, err := client.Login(context.Background(), &chat.LoginRequest{Name: username})
	if err != nil {
		msg := fmt.Sprintf("Login error: %s\n", err)
		g.Msgbox("Login error", msg)
		log.Print(msg)
		return
	}
	if resp.Status {
		msg := fmt.Sprintf("Successful login [id %d], message from server: %s\n", resp.Id, resp.Message)
		userid = resp.Id
		log.Print(msg)
		g.Msgbox("Login success", msg)
		time.Sleep(2 * time.Second)
		startChat()
	} else {
		msg := fmt.Sprintf("Rejected login: %s\n", resp.Message)
		g.Msgbox("Login error", msg)
		log.Print(msg)
	}
}

func loop() {
	if !logged {
		g.Window("Client login").Size(400, 200).Flags(g.WindowFlags(g.MasterWindowFlagsNotResizable)).Layout(
			g.Row(
				g.Label("User name"),
				g.InputText(&username),
			),
			g.Row(
				g.Label("Address"),
				g.InputText(&serverAddress),
			),
			g.Button("Login").OnClick(login),
			g.PrepareMsgbox(),
		)
	} else {
		g.SingleWindowWithMenuBar().Layout(
			g.SplitLayout(g.DirectionHorizontal, true, 200, g.Layout{
				g.Label("Chat participants"),
				g.RangeBuilder("Clients", clients, func(i int, v interface{}) g.Widget {
					str := v.(string)
					return g.Label(str)
				}),
			}, g.SplitLayout(g.DirectionVertical, true, -200, g.Layout{}, g.Layout{})),
		)
	}
}

func main() {
	logged = false
	wnd = g.NewMasterWindow("Gochat client", 400, 200, g.MasterWindowFlagsMaximized)
	wnd.Run(loop)
}

func startChat() {
	logged = true
	/*title := fmt.Sprintf("Awesome gRPC chat [%s]", username)
	w := g.NewMasterWindow(title, 800, 600, g.MasterWindowFlagsMaximized)
	w.Run(chatLoop)
	wnd = w*/
	/*w := (*a).NewWindow(title)
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
	})*/

	c := &client

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

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				log.Println("Connection closed by server")
			}
			if err != nil {
				log.Fatalf("Error receiving message: %v\n", err)
			}

			switch msg.Type {
			case "CLIENTS":
				log.Printf("Clients message received: %s\n", msg.Content)
				clientsList := strings.Split(msg.Content, ",")
				clients = make([]interface{}, 0)
				for _, cl := range clientsList {
					clients = append(clients, cl)
				}

				if len(msg.Name) > 0 {
					//byeStr := fmt.Sprintf("---- %s has left the chat ----", msg.Name)
					//chatWindow.Append(widget.NewLabel(byeStr))
				}
			case "MESSAGE":
				//messageStr := fmt.Sprintf("%s ---> %s", msg.Name, msg.Content)
				//chatWindow.Append(widget.NewLabel(messageStr))
				//chatWindowScroll.ScrollToBottom()
			default:
				log.Printf("Unknown message type %s: %v\n", msg.Type, msg)
			}
		}
	}()

	/*chatSend := widget.NewButton("Send", func() {
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
	w.Show()*/
}
