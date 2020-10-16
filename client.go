package main

import (
	"context"
	"fmt"
	"log"
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
	w := (*a).NewWindow("Awesome gRPC chat")
	w.CenterOnScreen()
	w.Resize(fyne.NewSize(800, 600))

	names := widget.NewVBox()
	names.Append(widget.NewLabel(username))
	nameScroll := widget.NewVScrollContainer(names)

	chatWindow := widget.NewVBox()
	chatWindowScroll := widget.NewVScrollContainer(chatWindow)

	chatMessage := widget.NewMultiLineEntry()
	chatMessage.Resize(fyne.NewSize(600, 200))
	chatSend := widget.NewButton("Send", func() {
		fmt.Println("We will send some cool message to the server: %s", chatMessage.Text)
	})

	chatBottom := widget.NewHBox(chatMessage, chatSend)

	chatArea := widget.NewVSplitContainer(chatWindowScroll, chatBottom)
	chatArea.SetOffset(10)

	genContainer := widget.NewHSplitContainer(nameScroll, chatArea)
	genContainer.SetOffset(-10)

	w.SetContent(genContainer)
	w.Show()
}
