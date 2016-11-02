package htmlquery

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/antchfx/gxpath"
	"github.com/antchfx/gxpath/xpath"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/htmlindex"
)

// CreateXPathNavigator creates a new xpath.NodeNavigator for the specified html.Node.
func CreateXPathNavigator(top *html.Node) xpath.NodeNavigator {
	return &htmlNodeNavigator{curr: top, root: top, attr: -1}
}

// Find searches the html.Node that matches by the specified XPath expr.
func Find(top *html.Node, expr string) []*html.Node {
	var elems []*html.Node
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	for t.MoveNext() {
		elems = append(elems, (t.Current().(*htmlNodeNavigator)).curr)
	}
	return elems
}

// FindOne searches the html.Node that matches by the specified XPath expr,
// and returns first element of matched html.Node.
func FindOne(top *html.Node, expr string) *html.Node {
	var elem *html.Node
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	if t.MoveNext() {
		elem = (t.Current().(*htmlNodeNavigator)).curr
	}
	return elem
}

// FindEach searches the html.Node and calls functions cb.
func FindEach(top *html.Node, expr string, cb func(int, *html.Node)) {
	t := gxpath.Select(CreateXPathNavigator(top), expr)
	i := 0
	for t.MoveNext() {
		cb(i, (t.Current().(*htmlNodeNavigator)).curr)
		i++
	}
}

func Load(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	charset := "utf-8"
	contyp := strings.ToLower(resp.Header.Get("Content-Type"))
	if i := strings.Index(contyp, "charset="); i > 0 {
		charset = contyp[i+8:]
	}
	e, err := htmlindex.Get(charset)
	if err != nil {
		return nil, err
	}
	if name, _ := htmlindex.Name(e); name != "utf8" {
		body := e.NewDecoder().Reader(resp.Body)
		return html.Parse(body)
	}
	return html.Parse(resp.Body)
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
	return strings.TrimSpace(buf.String())
}

func isSelfClosingTag(t string) bool {
	switch t {
	case "area", "hr", "img", "meta", "source", "br", "input":
		return true
	default:
		return false
	}
}

func outputXML(buf *bytes.Buffer, n *html.Node) {
	if n.Type == html.TextNode || n.Type == html.CommentNode {
		buf.WriteString(strings.TrimSpace(n.Data))
		return
	}

	buf.WriteString("<" + n.Data)
	for _, attr := range n.Attr {
		buf.WriteString(fmt.Sprintf(` %s="%s"`, attr.Key, attr.Val))
	}
	selfClosing := isSelfClosingTag(n.Data)
	if selfClosing {
		buf.WriteString("/>")
	} else {
		buf.WriteString(">")
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		outputXML(buf, child)
	}
	if !selfClosing {
		buf.WriteString(fmt.Sprintf("</%s>", n.Data))
	}
}

// OutputHTML returns the text including tags name.
func OutputHTML(n *html.Node) string {
	var buf bytes.Buffer
	outputXML(&buf, n)
	return buf.String()
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
	case html.DoctypeNode:
		// ignored <!DOCTYPE HTML> declare and as Root-Node type.
		return xpath.RootNode
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
