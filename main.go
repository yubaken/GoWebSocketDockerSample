package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/config"
	"github.com/kataras/iris/websocket"
	"os"
)

// In your main() function
type clientPage struct {
	Title string
	Host  string
}

func main() {
	api := iris.New()

	api.Static("/js", "./static/js", 1)

	api.Get("/", func(ctx *iris.Context) {
		ctx.Render("client.html", clientPage{"Client Page", ctx.HostString()})
	})

	// important staff
	websocketConfig := config.DefaultWebsocket()
	websocketConfig.Endpoint = "/my_endpoint" // the path which the websocket client should listen/registed to
	w := websocket.New(api, websocketConfig, api.Logger)
	// you created a new websocket server, you can create more than one... I leave that to you: w2:= websocket.New...; w2.OnConnection(...)
	// for default 'iris.' station use that: w := websocket.New(iris.DefaultIris, "/my_endpoint")
	var myChatRoom = "room1"
	w.OnConnection(func(c websocket.Connection) {

		c.Join(myChatRoom)

		c.On("chat", func(message string) {
			// to all except this connection ->
			//c.To(websocket.Broadcast).Emit("chat", "Message from: "+c.ID()+"-> "+message)

			// to the client ->
			//c.Emit("chat", "Message from myself: "+message)

			//send the message to the whole room,
			//all connections are inside this room will receive this message
			c.To(myChatRoom).Emit("chat", "From: "+c.ID()+": "+message)
		})

		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})
	})

	//

	api.Listen("0.0.0.0"+":80")
}