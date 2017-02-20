package xmlquery

import (
	"fmt"
	"strings"

	"github.com/antchfx/xpath"
)

// SelectElements finds child elements with the specified name.
func (n *Node) SelectElements(name string) []*Node {
	return Find(n, name)
}

// SelectElements finds child elements with the specified name.
func (n *Node) SelectElement(name string) *Node {
	return FindOne(n, name)
}

// SelectAttr returns the attribute value with the specified name.
func (n *Node) SelectAttr(name string) string {
	var local, space string
	local = name
	if i := strings.Index(name, ":"); i > 0 {
		space = name[:i]
		local = name[i+1:]
	}
	for _, attr := range n.Attr {
		if attr.Name.Local == local && attr.Name.Space == space {
			return attr.Value
		}
	}
	return ""
}

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func CreateXPathNavigator(top *Node) xpath.NodeNavigator {
	return &xmlNodeNavigator{curr: top, root: top, attr: -1}
}

func Find(top *Node, expr string) []*Node {
	t := xpath.Select(CreateXPathNavigator(top), expr)
	var elems []*Node
	for t.MoveNext() {
		elems = append(elems, (t.Current().(*xmlNodeNavigator)).curr)
	}
	return elems
}

func FindOne(top *Node, expr string) *Node {
	t := xpath.Select(CreateXPathNavigator(top), expr)
	var elem *Node
	if t.MoveNext() {
		elem = (t.Current().(*xmlNodeNavigator)).curr
	}
	return elem
}

// FindEach searches the html.Node and calls functions cb.
func FindEach(top *Node, expr string, cb func(int, *Node)) {
	t := xpath.Select(CreateXPathNavigator(top), expr)
	var i int
	for t.MoveNext() {
		cb(i, (t.Current().(*xmlNodeNavigator)).curr)
		i++
	}
}

type xmlNodeNavigator struct {
	root, curr *Node
	attr       int
}

func (x *xmlNodeNavigator) NodeType() xpath.NodeType {
	switch x.curr.Type {
	case CommentNode:
		return xpath.CommentNode
	case TextNode:
		return xpath.TextNode
	case DeclarationNode, DocumentNode:
		return xpath.RootNode
	case ElementNode:
		if x.attr != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	}
	panic(fmt.Sprintf("unknown XML node type: %v", x.curr.Type))
}

func (x *xmlNodeNavigator) LocalName() string {
	if x.attr != -1 {
		return x.curr.Attr[x.attr].Name.Local
	}
	return x.curr.Data

}

func (x *xmlNodeNavigator) Prefix() string {
	return x.curr.Namespace
}

func (x *xmlNodeNavigator) Value() string {
	switch x.curr.Type {
	case CommentNode:
		return x.curr.Data
	case ElementNode:
		if x.attr != -1 {
			return x.curr.Attr[x.attr].Value
		}
		return x.curr.InnerText()
	case TextNode:
		return x.curr.Data
	}
	return ""
}

func (x *xmlNodeNavigator) Copy() xpath.NodeNavigator {
	n := *x
	return &n
}

func (x *xmlNodeNavigator) MoveToRoot() {
	x.curr = x.root
}

func (x *xmlNodeNavigator) MoveToParent() bool {
	if node := x.curr.Parent; node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *xmlNodeNavigator) MoveToNextAttribute() bool {
	if x.attr >= len(x.curr.Attr)-1 {
		return false
	}
	x.attr++
	return true
}

func (x *xmlNodeNavigator) MoveToChild() bool {
	if node := x.curr.FirstChild; node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *xmlNodeNavigator) MoveToFirst() bool {
	if x.curr.PrevSibling == nil {
		return false
	}
	for {
		node := x.curr.PrevSibling
		if node == nil {
			break
		}
		x.curr = node
	}
	return true
}

func (x *xmlNodeNavigator) String() string {
	return x.Value()
}

func (x *xmlNodeNavigator) MoveToNext() bool {
	if node := x.curr.NextSibling; node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *xmlNodeNavigator) MoveToPrevious() bool {
	if node := x.curr.PrevSibling; node != nil {
		x.curr = node
		return true
	}
	return false
}

func (x *xmlNodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*xmlNodeNavigator)
	if !ok || node.root != x.root {
		return false
	}

	x.curr = node.curr
	x.attr = node.attr
	return true
}
