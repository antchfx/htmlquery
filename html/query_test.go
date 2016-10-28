package htmlquery

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

var doc = loadHtml()

func TestXPathSelect(t *testing.T) {
	if node := FindOne(doc, "/html/head/title"); node == nil {
		t.Fatal("cannot found any node")
	}
	if node := FindOne(doc, "//body[@bgcolor]"); node.Attr[0].Val != "ffffff" {
		t.Fatal("body bgcolor is not #ffffff")
	}
	if list := Find(doc, "//a"); len(list) != 2 {
		t.Fatal("count(//a)!=2")
	}
	if list := Find(doc, "//body/child::*"); len(list) != 9 { // ignored textnode
		t.Fatal("count(//body/child::*)!=9")
	}
}

func TestInnerText(t *testing.T) {
	title := FindOne(doc, "//title")
	if txt := InnerText(title); strings.TrimSpace(txt) != "your title here" {
		t.Fatalf("InnerText(//title): %s !=your title here", txt)
	}
	head := FindOne(doc, "/html/head")
	if txt := InnerText(head); strings.TrimSpace(txt) != "your title here" {
		t.Fatalf("InnerText(/html/head): %s !=your title here", txt)
	}
	img := FindOne(doc, "//img")
	if OutputHTML(img) != `<img src="clouds.jpg" align="bottom"/>` {
		t.Fatal(`OutputHTML(img)!='<img src="clouds.jpg" align="bottom"/>'`)
	}
}

func loadHtml() *html.Node {
	// http://help.websiteos.com/websiteos/example_of_a_simple_html_page.htm
	var str = `<!DOCTYPE html><html>
<head>
<title>your title here</title>
</head>
<body bgcolor="ffffff">
<center><img src="clouds.jpg" align="bottom"> </center>
<hr>
<a href="http://somegreatsite.com">link name</a>
is a link to another nifty site
<h1>this is a header</h1>
<h2>this is a medium header</h2>
send me mail at <a href="mailto:support@yourcompany.com">support@yourcompany.com</a>.
<p> this is a new paragraph!
<p> <b>this is a new paragraph!</b>
<br> <b><i>this is a new sentence without a paragraph break, in bold italics.</i></b>
<hr>
</body>
</html`
	node, err := html.Parse(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	return node
}
