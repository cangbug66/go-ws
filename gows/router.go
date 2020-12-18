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

func parsePattern(pattern string) []string {
	vs:=strings.Split(pattern,"/")

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

func (r *router) addRoute(pattern string,handler HandlerFunc)  {
	parts:=parsePattern(pattern)
	log.Printf("Route %s", pattern)
	fmt.Println("parts",parts)
	r.root.insert(pattern,parts,0)
	r.handlers[pattern] = handler
}

func (r *router) getRoute(pattern string) (*node,map[string]string) {
	searchParts := parsePattern(pattern)
	params := map[string]string{}
	n:=r.root.search(searchParts,0)
	if n!=nil{
		parts:=parsePattern(n.pattern)
		for index,part:=range parts{
			if part[0] == ':'{
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1{
				params[part[1:]] = strings.Join(searchParts[index:],"/")
				break
			}
		}
		return n, params
	}
	return nil,nil
}

func (r *router) handle(engine *Engine,c *Context) {

	Data:=c.Message.Data
	form:=&WsForm{}
	if err:=GetForm(Data,form);err!=nil{
		log.Println("json解析失败,参数:", string(Data))
		return
	}
	action := form.Action
	if action == ""{
		c.String("action is null")
		return
	}
	n,params:=r.getRoute(action)

	var middlewares []HandlerFunc
	for _,group:=range engine.groups{
		if strings.HasPrefix(action,group.prefix){
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c.handlers = middlewares
	if n!=nil{
		c.Params = params
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