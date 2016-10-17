package html

import (
	"bytes"
	"fmt"

	"github.com/antchfx/gxpath"
	"github.com/antchfx/gxpath/xpath"
	"golang.org/x/net/html"
)

// Selector is a XPath selector for HTML.
type Selector struct{}

func (s *Selector) Find(top *html.Node, expr string) []*html.Node {
	nav := &htmlNodeNavigator{curr: top, root: top, attr: -1}
	var elems []*html.Node
	t := s.Select(nav, expr)
	for t.MoveNext() {
		elems = append(elems, (t.Current().(*htmlNodeNavigator)).curr)
	}
	return elems
}

func (s *Selector) FindOne(top *html.Node, expr string) *html.Node {
	nav := &htmlNodeNavigator{curr: top, root: top, attr: -1}
	t := s.Select(nav, expr)
	var elem *html.Node
	if t.MoveNext() {
		elem = (t.Current().(*htmlNodeNavigator)).curr
	}
	return elem
}

func (s *Selector) Select(root xpath.NodeNavigator, expr string) *gxpath.NodeIterator {
	return gxpath.Select(root, expr)
}

type htmlNodeNavigator struct {
	root, curr *html.Node
	attr       int
}

func (h *htmlNodeNavigator) NodeType() xpath.NodeType {
	switch h.curr.Type {
	case html.CommentNode:
		return xpath.CommentNode
	case html.TextNode:
		return xpath.TextNode
	case html.DocumentNode:
		return xpath.RootNode
	case html.ElementNode:
		if h.attr != -1 {
			return xpath.AttributeNode
		}
		return xpath.ElementNode
	}
	panic(fmt.Sprintf("unknown HTML node type: %v", h.curr.Type))
}

func (h *htmlNodeNavigator) LocalName() string {
	if h.attr != -1 {
		return h.curr.Attr[h.attr].Key
	}
	return h.curr.Data
}

func (*htmlNodeNavigator) Prefix() string {
	return ""
}

func (h *htmlNodeNavigator) Value() string {
	switch h.curr.Type {
	case html.CommentNode:
		return h.curr.Data
	case html.ElementNode:
		if h.attr != -1 {
			return h.curr.Attr[h.attr].Val
		}
		return InnerText(h.curr)
	case html.TextNode:
		return h.curr.Data
	}
	return ""
}

func (h *htmlNodeNavigator) Copy() xpath.NodeNavigator {
	n := *h
	return &n
}

func (h *htmlNodeNavigator) MoveToRoot() {
	h.curr = h.root
}

func (h *htmlNodeNavigator) MoveToParent() bool {
	if node := h.curr.Parent; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *htmlNodeNavigator) MoveToNextAttribute() bool {
	if h.attr >= len(h.curr.Attr)-1 {
		return false
	}
	h.attr++
	return true
}

func (h *htmlNodeNavigator) MoveToChild() bool {
	if node := h.curr.FirstChild; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *htmlNodeNavigator) MoveToFirst() bool {
	if h.curr.PrevSibling == nil {
		return false
	}
	for {
		node := h.curr.PrevSibling
		if node == nil {
			break
		}
		h.curr = node
	}
	return true
}

func (h *htmlNodeNavigator) String() string {
	return h.Value()
}

func (h *htmlNodeNavigator) MoveToNext() bool {
	if node := h.curr.NextSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *htmlNodeNavigator) MoveToPrevious() bool {
	if node := h.curr.PrevSibling; node != nil {
		h.curr = node
		return true
	}
	return false
}

func (h *htmlNodeNavigator) MoveTo(other xpath.NodeNavigator) bool {
	node, ok := other.(*htmlNodeNavigator)
	if !ok || node.root != h.root {
		return false
	}

	h.curr = node.curr
	h.attr = node.attr
	return true
}

// InnerText returns the text between the start and end tags of the object.
func InnerText(n *html.Node) string {
	if n.Type == html.TextNode || n.Type == html.CommentNode {
		return n.Data
	}
	var buf bytes.Buffer
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		buf.WriteString(InnerText(child))
	}
	return buf.String()
}
