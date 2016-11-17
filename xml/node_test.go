package xmlquery

import (
	"strings"
	"testing"
)

func findNode(root *Node, name string) *Node {
	node := root.FirstChild
	for {
		if node == nil || node.Data == name {
			break
		}
		node = node.NextSibling
	}
	return node
}

func childNodes(root *Node, name string) []*Node {
	var list []*Node
	node := root.FirstChild
	for {
		if node == nil {
			break
		}
		if node.Data == name {
			list = append(list, node)
		}
		node = node.NextSibling
	}
	return list
}

func testNode(t *testing.T, n *Node, expected string) {
	if n.Data != expected {
		t.Fatalf("expected node name is %s,but got %s", expected, n.Data)
	}
}

func testAttr(t *testing.T, n *Node, name, expected string) {
	for _, attr := range n.Attr {
		if attr.Name.Local == name && attr.Value == expected {
			return
		}
	}
	t.Fatalf("not found attribute %s in the node %s", name, n.Data)
}

func testValue(t *testing.T, val, expected string) {
	if val != expected {
		t.Fatalf("expected value is %s,but got %s", expected, val)
	}
}

func TestParse(t *testing.T) {
	s := `<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
<book>
  <title lang="en">Harry Potter</title>
  <price>29.99</price>
</book>
<book>
  <title lang="en">Learning XML</title>
  <price>39.95</price>
</book>
</bookstore>`
	root, err := ParseXML(strings.NewReader(s))
	if err != nil {
		t.Error(err)
	}
	if root.Type != DocumentNode {
		t.Fatal("top node of tree is not DocumentNode")
	}

	declarNode := root.FirstChild
	if declarNode.Type != DeclarationNode {
		t.Fatal("first child node of tree is not DeclarationNode")
	}

	if declarNode.Attr[0].Name.Local != "version" && declarNode.Attr[0].Value != "1.0" {
		t.Fatal("version attribute not expected")
	}

	bookstore := root.LastChild
	if bookstore.Data != "bookstore" {
		t.Fatal("bookstore elem not found")
	}
	if bookstore.FirstChild.Data != "\n" {
		t.Fatal("first child node of bookstore is not empty node(\n)")
	}
	books := childNodes(bookstore, "book")
	if len(books) != 2 {
		t.Fatalf("expected book element count is 2, but got %d", len(books))
	}
	// first book element
	testNode(t, findNode(books[0], "title"), "title")
	testAttr(t, findNode(books[0], "title"), "lang", "en")
	testValue(t, findNode(books[0], "price").InnerText(), "29.99")
	testValue(t, findNode(books[0], "title").InnerText(), "Harry Potter")

	// second book element
	testNode(t, findNode(books[1], "title"), "title")
	testAttr(t, findNode(books[1], "title"), "lang", "en")
	testValue(t, findNode(books[1], "price").InnerText(), "39.95")

	testValue(t, books[0].OutputXML(), `<book><title lang="en">Harry Potter</title><price>29.99</price></book>`)
}

func TestTooNested(t *testing.T) {
	s := `<?xml version="1.0" encoding="UTF-8"?>
    <AAA> 
        <BBB> 
            <DDD> 
                <CCC> 
                    <DDD/> 
                    <EEE/> 
                </CCC> 
            </DDD> 
        </BBB> 
        <CCC> 
            <DDD> 
                <EEE> 
                    <DDD> 
                        <FFF/> 
                    </DDD> 
                </EEE> 
            </DDD> 
        </CCC> 
     </AAA>`
	root, err := ParseXML(strings.NewReader(s))
	if err != nil {
		t.Error(err)
	}
	aaa := findNode(root, "AAA")
	if aaa == nil {
		t.Fatal("AAA node not exists")
	}
	ccc := aaa.LastChild
	if ccc.Data != "CCC" {
		t.Fatalf("expected node is CCC,but got %s", ccc.Data)
	}
	bbb := ccc.PrevSibling
	if bbb.Data != "BBB" {
		t.Fatalf("expected node is bbb,but got %s", bbb.Data)
	}
	ddd := findNode(bbb, "DDD")
	testNode(t, ddd, "DDD")
	testNode(t, ddd.LastChild, "CCC")
}
