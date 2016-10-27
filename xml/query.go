package xmlquery

import (
	"fmt"

	"github.com/antchfx/gxpath"
	"github.com/antchfx/gxpath/xpath"
	"github.com/antchfx/xml"
)

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func CreateXPathNavigator(top *xml.Node) xpath.NodeNavigator {
	return &xmlNodeNavigator{curr: top, root: top, attr: -1}
}

func Find(top *xml.Node, expr string) []*xml.Node {
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	var elems []*xml.Node
	for t.MoveNext() {
		elems = append(elems, (t.Current().(*xmlNodeNavigator)).curr)
	}
	return elems
}

func FindOne(top *xml.Node, expr string) *xml.Node {
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	var elem *xml.Node
	if t.MoveNext() {
		elem = (t.Current().(*xmlNodeNavigator)).curr
	}
	return elem
}

// FindEach searches the html.Node and calls functions cb.
func FindEach(top *xml.Node, expr string, cb func(int, *xml.Node)) {
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	i := 0
	if t.MoveNext() {
		cb(i, (t.Current().(*xmlNodeNavigator)).curr)
		i++
	}
}

type xmlNodeNavigator struct {
	root, curr *xml.Node
	attr       int
}

func (x *xmlNodeNavigator) NodeType() xpath.NodeType {
	switch x.curr.Type {
	case xml.CommentNode:
		return xpath.CommentNode
	case xml.TextNode:
		return xpath.TextNode
	case xml.DeclarationNode, xml.DocumentNode:
		return xpath.RootNode
	case xml.ElementNode:
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
	case xml.CommentNode:
		return x.curr.Data
	case xml.ElementNode:
		if x.attr != -1 {
			return x.curr.Attr[x.attr].Value
		}
		return x.curr.InnerText()
	case xml.TextNode:
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
