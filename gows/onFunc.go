package gows

type OnFunc struct {
	OnMessageFunc HandlerFunc
	OnCloseFunc   HandlerFunc
	OnOpenFunc    HandlerFunc
	OnErrorFunc   HandlerFunc
}

func (engine *Engine) OnOpenFunc(f HandlerFunc) *Engine {
	engine.OnFunc.OnOpenFunc = f
	return engine
}
func (engine *Engine) OnMessageFunc(f HandlerFunc) *Engine {
	engine.OnFunc.OnMessageFunc = f
	return engine
}
func (engine *Engine) OnCloseFunc(f HandlerFunc) *Engine {
	engine.OnFunc.OnCloseFunc = f
	return engine
}
func (engine *Engine) OnErrorFunc(f HandlerFunc) *Engine {
	engine.OnFunc.OnErrorFunc = f
	return engine
}