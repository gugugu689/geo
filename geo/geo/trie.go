package geo

import "strings"

//node 树的节点
type node struct {
	//pattern 待匹配路由
	pattern string
	//part	节点参数
	part string
	//wildChild 是否模糊匹配 当path含有 * ： 为true
	wildChild bool
	//children 子节点切片
	children []*node
	//fullPath 完整路径
	fullPath string
}

//methodTree 方法树 一个方法对应一个树
type methodTree struct {
	method string
	root   *node
}

//methodTrees 方法树切片，对应http方法
type methodTrees []methodTree

func (trees methodTrees) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	//root 方法树根节点
	root := trees.get(method)

	if root == nil {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]

			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

//get 根据请求方法找对应的方法树根节点
func (trees methodTrees) get(method string) *node {
	for _, tree := range trees {
		if tree.method == method {
			return tree.root
		}
	}
	return nil
}

//parsePattern 将路由按照 '/' 拆分
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0) //make一个大小为 0 的 string 类型数组
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

//addRoute 添加路由节点
func (n *node) addRoute(pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	n.insert(pattern, parts, 0)
}

//insert 插入节点
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	child := n.matchInsert(part)
	if child == nil {
		child = &node{part: part, wildChild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//matchInsert 插入节点所用的搜索
func (n *node) matchInsert(part string) *node {
	for _, v := range n.children {
		if v.part == part || v.wildChild == true {
			return v
		}
	}
	return nil
}

//matchSearch 搜索节点所用的搜索
func (n *node) matchSearch(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.wildChild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

//search
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchSearch(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}
