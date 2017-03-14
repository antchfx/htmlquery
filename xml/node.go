package xmlquery

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/net/html/charset"
)

// A NodeType is the type of a Node.
type NodeType uint

const (
	DocumentNode NodeType = iota
	DeclarationNode
	ElementNode
	TextNode
	CommentNode
)

type Node struct {
	Parent, FirstChild, LastChild, PrevSibling, NextSibling *Node

	Type         NodeType
	Data         string
	Prefix       string
	NamespaceURI string
	Attr         []xml.Attr

	level int // node level in the tree
}

// InnerText returns the text between the start and end tags of the object.
func (n *Node) InnerText() string {
	var output func(*bytes.Buffer, *Node)
	output = func(buf *bytes.Buffer, n *Node) {
		switch n.Type {
		case TextNode:
			buf.WriteString(n.Data)
			return
		case CommentNode:
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			output(buf, child)
		}
	}

	var buf bytes.Buffer
	output(&buf, n)
	return buf.String()
}

func outputXML(buf *bytes.Buffer, n *Node) {
	if n.Type == TextNode || n.Type == CommentNode {
		buf.WriteString(strings.TrimSpace(n.Data))
		return
	}
	buf.WriteString("<" + n.Data)
	for _, attr := range n.Attr {
		if attr.Name.Space != "" {
			buf.WriteString(fmt.Sprintf(` %s:%s="%s"`, attr.Name.Space, attr.Name.Local, attr.Value))
		} else {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, attr.Name.Local, attr.Value))
		}
	}
	buf.WriteString(">")
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		outputXML(buf, child)
	}
	buf.WriteString(fmt.Sprintf("</%s>", n.Data))
}

// OutputXML returns the text that including tags name.
func (n *Node) OutputXML() string {
	var buf bytes.Buffer
	outputXML(&buf, n)
	return buf.String()
}

func addAttr(n *Node, key, val string) {
	var attr xml.Attr
	if i := strings.Index(key, ":"); i > 0 {
		attr = xml.Attr{
			Name:  xml.Name{Space: key[:i], Local: key[i+1:]},
			Value: val,
		}
	} else {
		attr = xml.Attr{
			Name:  xml.Name{Local: key},
			Value: val,
		}
	}

	n.Attr = append(n.Attr, attr)
}

func addChild(parent, n *Node) {
	n.Parent = parent
	if parent.FirstChild == nil {
		parent.FirstChild = n
	} else {
		parent.LastChild.NextSibling = n
		n.PrevSibling = parent.LastChild
	}

	parent.LastChild = n
}

func addSibling(sibling, n *Node) {
	n.Parent = sibling.Parent
	sibling.NextSibling = n
	n.PrevSibling = sibling
	if sibling.Parent != nil {
		sibling.Parent.LastChild = n
	}
}

// LoadURL loads the XML document from the specified URL.
func LoadURL(url string) (*Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}
	return ParseXML(r)
}

func parse(r io.Reader) (*Node, error) {
	var (
		decoder      = xml.NewDecoder(r)
		doc          = &Node{Type: DocumentNode}
		space2prefix = make(map[string]string)
		level        = 0
		declared     = false
	)
	var prev *Node = doc
	for {
		tok, err := decoder.Token()
		switch {
		case err == io.EOF:
			goto quit
		case err != nil:
			return nil, err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			if !declared {
				return nil, errors.New("xml: document is invalid")
			}
			node := &Node{
				Type:         ElementNode,
				Data:         tok.Name.Local,
				Prefix:       space2prefix[tok.Name.Space],
				NamespaceURI: tok.Name.Space,
				Attr:         tok.Attr,
				level:        level,
			}
			for _, att := range tok.Attr {
				if att.Name.Space == "xmlns" {
					space2prefix[att.Value] = att.Name.Local
				}
			}
			//fmt.Println(fmt.Sprintf("start > %s : %d", node.Data, level))
			if level == prev.level {
				addSibling(prev, node)
			} else if level > prev.level {
				addChild(prev, node)
			} else if level < prev.level {
				for i := prev.level - level; i > 1; i-- {
					prev = prev.Parent
				}
				addSibling(prev.Parent, node)
			}
			prev = node
			level++
		case xml.EndElement:
			level--
		case xml.CharData:
			node := &Node{Type: TextNode, Data: string(tok), level: level}
			if level == prev.level {
				addSibling(prev, node)
			} else if level > prev.level {
				addChild(prev, node)
			}
		case xml.Comment:
			node := &Node{Type: CommentNode, Data: string(tok), level: level}
			if level == prev.level {
				addSibling(prev, node)
			} else if level > prev.level {
				addChild(prev, node)
			}
		case xml.ProcInst: // Processing Instruction
			if declared || (!declared && tok.Target != "xml") {
				return nil, errors.New("xml: document is invalid")
			}
			level++
			node := &Node{Type: DeclarationNode, level: level}
			pairs := strings.Split(string(tok.Inst), " ")
			for _, pair := range pairs {
				pair = strings.TrimSpace(pair)
				if i := strings.Index(pair, "="); i > 0 {
					addAttr(node, pair[:i], strings.Trim(pair[i+1:], `"`))
				}
			}
			declared = true
			if level == prev.level {
				addSibling(prev, node)
			} else if level > prev.level {
				addChild(prev, node)
			}
			prev = node
		case xml.Directive:
		}

	}
quit:
	return doc, nil
}

// Parse returns the parse tree for the XML from the given Reader.
func Parse(r io.Reader) (*Node, error) {
	return parse(r)
}

// Deprecated,Parse returns the parse tree for the XML from the given Reader.
func ParseXML(r io.Reader) (*Node, error) {
	return parse(r)
}
