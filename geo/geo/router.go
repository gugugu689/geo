package geo

import (
	"net/http"
	"path"
)

//Router 所有路由映射和路由组
type Router interface {
	Routes
	Group(string, []HandlerFunc) *RouterGroup
}

//Routes 定义所有的路由映射
type Routes interface {
	Use([]HandlerFunc) Routes

	GET(string, HandlerFunc) Routes
	POST(string, HandlerFunc) Routes
	DELETE(string, HandlerFunc) Routes
	PUT(string, HandlerFunc) Routes
}

//RouterGroup 路由组
type RouterGroup struct {
	//Handlers 中间件
	Handlers []HandlerFunc
	//basePath 前缀路径
	basePath string
	engine   *Engine
	root     bool
}

//Group 创建路由组方法
func (group *RouterGroup) Group(relativePath string, handlers []HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers: group.combineHandlers(handlers),
		basePath: group.calculateAbsolutePath(relativePath),
		engine:   group.engine,
	}
}

//Use 中间件使用
func (group *RouterGroup) Use(middleware []HandlerFunc) Routes {
	group.Handlers = append(group.Handlers, middleware...)
	return group.returnSelf()
}

//handle 处理请求方法
func (group *RouterGroup) handle(httpMethod, relativePath string, handler HandlerFunc) Routes {
	absolutePath := group.calculateAbsolutePath(relativePath)
	group.engine.addRoute(httpMethod, absolutePath, handler)
	return group.returnSelf()
}
func (group *RouterGroup) POST(relativePath string, handler HandlerFunc) Routes {
	return group.handle(http.MethodPost, relativePath, handler)
}
func (group *RouterGroup) GET(relativePath string, handler HandlerFunc) Routes {
	return group.handle(http.MethodGet, relativePath, handler)
}
func (group *RouterGroup) DELETE(relativePath string, handler HandlerFunc) Routes {
	return group.handle(http.MethodDelete, relativePath, handler)
}
func (group *RouterGroup) PUT(relativePath string, handler HandlerFunc) Routes {
	return group.handle(http.MethodPut, relativePath, handler)
}

//returnSelf 返回engine或者路由俎
func (group *RouterGroup) returnSelf() Routes {
	if group.root {
		return group.engine
	}
	return group
}

//calculateAbsolutePath 计算绝对路径
func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if group.basePath == "" {
		return group.basePath
	}
	finalPath := path.Join(group.basePath, relativePath)
	if relativePath[len(relativePath)-1] == '/' && finalPath[len(finalPath)-1] != '/' {
		return finalPath + "/"
	}
	return finalPath
}

//combineHandlers 组合合并handlers
func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, finalSize)
	copy(mergedHandlers, group.Handlers)
	copy(mergedHandlers[len(group.Handlers):], handlers)
	return mergedHandlers
}
