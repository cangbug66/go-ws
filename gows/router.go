package gows

import (
	"fmt"
	"log"
	"strings"
)

type router struct {
	root *node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		root: &node{},
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePath(path string) []string {
	vs:=strings.Split(path,"/")

	parts := []string{}
	for _,item:=range vs{
		if item != ""{
			parts = append(parts,item)
			if(item[0] == '*'){
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(path string,handler HandlerFunc)  {
	parts:=parsePath(path)
	log.Printf("Route: %s", path)
	r.root.insert(path,parts,0)
	r.handlers[path] = handler
}

func (r *router) getRoute(path string) *node {
	searchParts := parsePath(path)
	n:=r.root.search(searchParts,0)
	if n!=nil{
		return n
	}
	return nil
}

func (r *router) handle(engine *Engine,c *Context) {

	Data:=c.Message.Data
	form:=&WsForm{}
	if err:=GetForm(Data,form);err!=nil{
		log.Println("json解析失败:", string(Data))
		return
	}
	action := form.Action
	if action == ""{
		c.String("action is null")
		return
	}
	n:=r.getRoute(action)

	var middlewares []HandlerFunc
	for _,group:=range engine.groups{
		if strings.HasPrefix(action,group.prefix){
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.handlers = middlewares
	if n!=nil{
		if hanlder,ok := r.handlers[action];ok{
			c.handlers = append(c.handlers,hanlder)
		}
	}else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.Conn.WriteMessage(1,[]byte(fmt.Sprintf("404 NOT FOUND: %s\n", action)))
		})
	}
	c.Action = form.Action
	c.Next()
}