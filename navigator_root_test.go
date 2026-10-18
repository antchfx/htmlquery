package htmlquery

import (
	"strings"
	"testing"

	"github.com/antchfx/xpath"
)

func TestMoveToRootFromAttribute(t *testing.T) {
	doc, err := Parse(strings.NewReader(`<div id="example"><span></span></div>`))
	if err != nil {
		t.Fatal(err)
	}
	iter := xpath.MustCompile("//div/@id").Select(CreateXPathNavigator(doc))
	if !iter.MoveNext() {
		t.Fatal("attribute not found")
	}
	nav := iter.Current().Copy()
	nav.MoveToRoot()
	if !nav.MoveToChild() {
		t.Fatal("root still has an attribute cursor and cannot move to its child")
	}
	if nav.NodeType() != xpath.ElementNode || nav.LocalName() != "html" {
		t.Fatalf("unexpected child after root navigation: %v", nav)
	}
}
