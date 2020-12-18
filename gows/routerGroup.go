package gows

type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	parent *RouterGroup
	engine *Engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine:=group.engine
	newGroup:=&RouterGroup{
		prefix:group.prefix+prefix,
		parent:group,
		engine:engine,
	}
	engine.groups = append(engine.groups,newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares  ...HandlerFunc){
	group.middlewares = append(group.middlewares,middlewares...)
}

func (group *RouterGroup) addRoute(pattern string, handler HandlerFunc) {
	pattern = group.prefix + pattern
	group.engine.router.addRoute(pattern,handler)
}

func (group *RouterGroup) Route(pattern string, handler HandlerFunc)  {
	group.addRoute(pattern,handler)
}