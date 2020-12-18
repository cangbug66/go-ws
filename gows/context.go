package gows

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Type int
	Data []byte
}

func NewMessage(t int, data []byte) *Message {
	return &Message{Type: t, Data: data}
}

type Context struct {
	Conn *websocket.Conn
	Writer http.ResponseWriter
	Req    *http.Request
	handlers []HandlerFunc
	index int
	Action string
	Params map[string]string
	Message *Message
}

func NewContext(conn *websocket.Conn,w http.ResponseWriter, r *http.Request,message *Message) *Context {
	return &Context{
		Conn:conn,
		Writer:w,
		Req:r,
		index:-1,
		Params: map[string]string{},
		Message:message,
		handlers:[]HandlerFunc{},
	}
}

func (c *Context) Next()  {
	c.index++
	s:=len(c.handlers)
	for ; c.index<s;c.index++{
		c.handlers[c.index](c)
	}
}

func (this *Context) Abort() {
	this.index = 9999
}

func (this *Context) AbortWithError()  {
	this.Abort()
}

func (c *Context) String(msg string) {
	if err:=c.Conn.WriteMessage(1, []byte(msg));err!=nil{
		log.Println("websocket send message err：",err)
	}
}