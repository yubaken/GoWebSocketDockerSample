package main

import (
	"fmt"

	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	//"database/sql"
	//"log"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	sb "github.com/dropbox/godropbox/database/sqlbuilder"
	"database/sql"
	"log"
	"github.com/spf13/viper"
	"bytes"
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

var db *sql.DB
var dbErr error

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("設定ファイル読み込みエラー: %s \n", err))
	}
	var dataSource bytes.Buffer
	dataSource.WriteString(viper.GetString("db.user"))
	dataSource.WriteString(":")
	dataSource.WriteString(viper.GetString("db.password"))
	dataSource.WriteString("@tcp(")
	dataSource.WriteString(viper.GetString("db.url"))
	dataSource.WriteString(":")
	dataSource.WriteString(viper.GetString("db.port"))
	dataSource.WriteString(")/chat")
	dataSource.WriteString("?interpolateParams=true&collation=utf8mb4_bin")

	db, dbErr = sql.Open("mysql",dataSource.String())
	if dbErr != nil {
		log.Fatal(dbErr)
	}

	iris.Static("/js", "./static/js", 1)

	iris.Get("/", func(ctx *iris.Context) {
		ctx.Render("client.html", clientPage{"Client Page", ctx.HostString()})
	})

	iris.Get("/sub", func(ctx *iris.Context) {
		ctx.Render("client2.html", clientPage{"Client Page", ctx.HostString()})
	})

	// important staff
	iris.Config.Websocket.Endpoint = "/my_endpoint"
	// you created a new websocket server, you can create more than one... I leave that to you: w2:= websocket.New...; w2.OnConnection(...)
	// for default 'iris.' station use that: w := websocket.New(iris.DefaultIris, "/my_endpoint")
	iris.Websocket.OnConnection(func(c websocket.Connection) {

		c.On("init", func(message string) {
			var d data
			json.Unmarshal([]byte(message), &d)
			c.Emit("chat", "join room: " + d.Room)
			if err != nil {
				log.Fatal(err)
			}
			t := sb.NewTable(
				"chats",
				sb.StrColumn("roomid",sb.UTF8,sb.UTF8CaseSensitive,false),
				sb.StrColumn("text",sb.UTF8,sb.UTF8CaseSensitive,false,),
			)
			query, _ := t.Select(t.C("roomid"), t.C("text")).String("chat")
			rows, err := db.Query(query)
			print(rows)
			if err != nil {
				log.Fatal(err)
			}
			for rows.Next() {
				var roomid string
				var text string
				if err := rows.Scan(&roomid, &text); err != nil {
					log.Fatal(err)
				}
				c.Emit("chat", text)
			}
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
			_, err = db.Exec(
				`INSERT INTO chats (roomid, text) VALUES (?, ?) `,
				d.Room,
				d.Msg,
			)
			if err != nil {
				log.Fatal(err)
			}
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