package main

import "github.com/kataras/iris"

func main() {
	iris.Get("/", func(ctx *iris.Context) {
		ctx.Write("Hi %s", "iris")
	})
	iris.Listen(":8080")
	//err := iris.ListenWithErr(":8080")
}