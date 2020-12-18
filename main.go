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
	a.Route("/b", func(context *gows.Context) {
		context.String("hello")
	})
	ws.Route("test", func(context *gows.Context) {
		context.String("test 测试")
	})

	ws.OnOpenFunc(func(context *gows.Context) {
		fmt.Println("open")
	})
	ws.OnCloseFunc(func(context *gows.Context) {
		fmt.Println("close")
	})
	ws.OnMessageFunc(func(context *gows.Context) {
		fmt.Println("message")
	})
	ws.Run(":8686","/ws")
}
