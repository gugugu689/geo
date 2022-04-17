package geo

import (
	"net/http"
	"strings"
)

//HandlerFunc 路由处理函数
type HandlerFunc func(c *Context)

type Engine struct {
	*RouterGroup
	groups []*RouterGroup
	//trees 方法树
	trees methodTrees
	//mappingHandlers 路由映射的处理方法
	mappingHandlers map[string]HandlerFunc
}

//addRoute 添加路由
func (engine *Engine) addRoute(httpMethod, absolutePath string, handlers HandlerFunc) {
	//找method对应的方法树 如果未空 则新建方法树根节点
	root := engine.trees.get(httpMethod)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		//将新的方法树添加进方法树切片
		engine.trees = append(engine.trees, methodTree{method: httpMethod, root: root})
	}
	root.addRoute(absolutePath, handlers)
	engine.mappingHandlers[absolutePath] = handlers
}

//New 新建Engine实例
func New() *Engine {
	engine := &Engine{
		RouterGroup: &RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		trees:           make([]methodTree, 0),
		mappingHandlers: map[string]HandlerFunc{},
	}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	engine.RouterGroup.engine = engine
	return engine
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.basePath) {
			middlewares = append(middlewares, group.Handlers...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	c.engine = engine
	engine.handleHTTP(c)
}
func (engine *Engine) handleHTTP(c *Context) {
	n, params := engine.trees.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, engine.mappingHandlers[key])
		//r.handlers[key](c)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
