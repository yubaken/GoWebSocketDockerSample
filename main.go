package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	//"database/sql"
	//"log"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
)

// In your main() function
type clientPage struct {
	Title string
	Host  string
}

type data struct {
	Room string
	Msg  string
}

func main() {
	iris.Static("/js", "./static/js", 1)

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("client.html", clientPage{"Client Page", ctx.HostString()})
	})

	iris.Get("/sub", func(ctx *iris.Context) {
		ctx.Render("client2.html", clientPage{"Client Page", ctx.HostString()})
	})

	//db, err := sql.Open("mysql", "root:password@db/chat?interpolateParams=true&collation=utf8mb4_bin")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//if _, err := db.Exec("SELECT SLEEP(?)", 42); err != nil {
	//	log.Fatal(err)
	//}
	// important staff
	iris.Config.Websocket.Endpoint = "/my_endpoint"
	// you created a new websocket server, you can create more than one... I leave that to you: w2:= websocket.New...; w2.OnConnection(...)
	// for default 'iris.' station use that: w := websocket.New(iris.DefaultIris, "/my_endpoint")
	iris.Websocket.OnConnection(func(c websocket.Connection) {

		c.On("init", func(message string) {
			var d data
			json.Unmarshal([]byte(message), &d)
			c.Emit("chat", "join room: " + d.Room)
			c.Join(d.Room)
		})

		c.On("chat", func(message string) {
			var d data
			json.Unmarshal([]byte(message), &d)

			// to all except this connection ->
			//c.To(websocket.Broadcast).Emit("chat", "Message from: "+c.ID()+"-> "+message)

			// to the client ->
			//c.Emit("chat", "Message from myself: "+message)

			//send the message to the whole room,
			//all connections are inside this room will receive this message
			c.To(d.Room).Emit("chat", "From: "+c.ID()+": "+d.Msg)
		})

		c.On("leave", func(message string) {
			var d data
			json.Unmarshal([]byte(message), &d)

			c.Leave(d.Room)
			c.To(d.Room).Emit("chat", "leave room: " + d.Room)
		})

		c.OnDisconnect(func() {
			fmt.Printf("\nConnection with ID: %s has been disconnected!", c.ID())
		})
	})

	iris.Listen("0.0.0.0:80")
}