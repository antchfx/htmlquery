package xquery

import (
	"fmt"

	"github.com/antchfx/gxpath"
	"github.com/antchfx/gxpath/xpath"
	"github.com/antchfx/xquery/html"
	"github.com/antchfx/xquery/xml"
)

// Type represents the documents types for XPath.
type Type int

const (
	HTML Type = iota
	XML
)

// Selector is an interface for XPath search.
type Selector interface {
	Select(xpath.NodeNavigator, string) *gxpath.NodeIterator
}

// New returns new Selector for the specified documents type.
func New(typ Type) Selector {
	switch typ {
	case HTML:
		return &html.Selector{}
	case XML:
		return &xml.Selector{}
	}
	panic(fmt.Errorf("unknown type : %d", typ))
}
