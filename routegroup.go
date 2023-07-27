package btcp

import "strconv"

type IRouter interface {
	IRoutes
	Group(string, ...HandlerFunc) *RouterGroup
}

type IRoutes interface {
	Use(handlers ...HandlerFunc) IRoutes
	Handle(code uint32, handlers ...HandlerFunc) IRoutes
}

type HandlerFunc func(request *Request)
type HandlersChain []HandlerFunc

type RouterGroup struct {
	Apis     map[uint32][]HandlerFunc
	Handlers []HandlerFunc
}

func NewRouteGroup() *RouterGroup {
	return &RouterGroup{
		Apis:     make(map[uint32][]HandlerFunc),
		Handlers: make([]HandlerFunc, 0),
	}
}

func (g *RouterGroup) Use(handlers ...HandlerFunc) IRoutes {
	g.Handlers = append(g.Handlers, handlers...)
	return g
}

func (g *RouterGroup) Handle(code uint32, handlers ...HandlerFunc) IRoutes {
	if _, ok := g.Apis[code]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(code)))
	}
	finalSize := len(g.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, g.Handlers)
	copy(mergedHandlers[len(g.Handlers):], handlers)
	g.Apis[code] = append(g.Apis[code], mergedHandlers...)
	return g
}
