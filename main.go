package main

import (
	"fmt"
	"gows"
)

func main() {
	ws:= gows.New()
	ws.Use(func(context *gows.Context) {
		fmt.Println("全局中间件")
	})
	a:=ws.Group("/a")
	a.Use(func(context *gows.Context) {
		fmt.Println("分组中间件")
	})
	a.Route("/b/:id", func(context *gows.Context) {
		fmt.Println("id",context.Params["id"])
		context.String("hello")
	})
	ws.Route("test", func(context *gows.Context) {
		context.String("test 测试")
	})
	ws.Run(":8686","/ws")
}
