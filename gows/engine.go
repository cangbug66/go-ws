package gows

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	upgrader    websocket.Upgrader
	router *router
	*RouterGroup
	groups []*RouterGroup
	*OnFunc
}

func New() *Engine {
	engine:=&Engine{
				router: newRouter(),
				upgrader:websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
				OnFunc:&OnFunc{},
			}
	engine.RouterGroup = &RouterGroup{engine:engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}

	return engine
}

func (engine *Engine) Run(addr string,path string) (err error) {
	log.Println(fmt.Sprintf("server listen addr: '%v' , path: '%v'", addr, path))
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		engine.wsServe(w, r)
	})
	return http.ListenAndServe(addr, nil)
}

func (engine *Engine) wsServe(w http.ResponseWriter, r *http.Request) {
	c, err := engine.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("http -> ws upgrade err:", err)
		return
	}
	go engine.ReadMessage(c,w,r)
}

func (engine *Engine) ReadMessage(conn *websocket.Conn,w http.ResponseWriter, r *http.Request) {
	defer conn.Close()
	conn.SetCloseHandler(func(code int, text string) error {
		if engine.OnFunc.OnCloseFunc != nil {
			context := NewContext(conn,w,r,NewMessage(0,[]byte{}))
			engine.OnFunc.OnCloseFunc(context)
		}
		return nil
	})
	if engine.OnFunc.OnOpenFunc != nil {
		context := NewContext(conn,w,r,NewMessage(0,[]byte{}))
		engine.OnFunc.OnOpenFunc(context)
	}
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read err:", err)
			break
		}
		log.Printf("recv: %s", message)
		context := NewContext(conn,w,r,NewMessage(messageType,message))
		if(engine.OnFunc.OnMessageFunc != nil){
			engine.OnFunc.OnMessageFunc(context)
			continue
		}
		engine.router.handle(engine,context)
	}
}
