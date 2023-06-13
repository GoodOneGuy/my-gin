package my_gin

import "strings"

type trieNode struct {
	pattern  string
	path     string
	children []*trieNode
	isWild   bool
}

func (n *trieNode) matchChild(path string) *trieNode {
	for _, child := range n.children {
		if child.path == path || child.isWild {
			return child
		}
	}

	return nil
}

func (n *trieNode) matchChildren(path string) []*trieNode {
	trieNodes := make([]*trieNode, 0)
	for _, child := range n.children {
		if child.path == path || child.isWild {
			trieNodes = append(trieNodes, child)
		}
	}

	return trieNodes
}

func (n *trieNode) insert(pattern string, paths []string, height int) {
	if len(paths) == height {
		n.pattern = pattern
		return
	}

	path := paths[height]
	child := n.matchChild(path)
	if child == nil {
		child = &trieNode{path: path, isWild: path[0] == ':' || path[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, paths, height+1)
}

func (n *trieNode) search(paths []string, height int) *trieNode {
	if len(paths) == height || strings.HasPrefix(n.path, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	path := paths[height]
	children := n.matchChildren(path)

	for _, child := range children {
		result := child.search(paths, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
