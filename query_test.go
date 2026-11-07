package htmlquery

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/antchfx/xpath"
	"golang.org/x/net/html"
)

const htmlSample = `<!DOCTYPE html><html lang="en-US">
<head>
<title>Hello,World!</title>
</head>
<body>
<div class="container">
<header>
	<!-- Logo -->
   <h1>City Gallery</h1>
</header>  
<nav>
  <ul>
    <li><a href="/London">London</a></li>
    <li><a href="/Paris">Paris</a></li>
    <li><a href="/Tokyo">Tokyo</a></li>
  </ul>
</nav>
<article>
  <h1>London</h1>
  <img src="pic_mountain.jpg" alt="Mountain View" style="width:304px;height:228px;">
  <p>London is the capital city of England. It is the most populous city in the  United Kingdom, with a metropolitan area of over 13 million inhabitants.</p>
  <p>Standing on the River Thames, London has been a major settlement for two millennia, its history going back to its founding by the Romans, who named it Londinium.</p>
</article>
<footer>Copyright &copy; W3Schools.com</footer>
</div>
</body>
</html>
`

var testDoc = loadHTML(htmlSample)

func BenchmarkSelectorCache(b *testing.B) {
	DisableSelectorCache = false
	for i := 0; i < b.N; i++ {
		getQuery("/AAA/BBB/DDD/CCC/EEE/ancestor::*")
	}
}

func BenchmarkDisableSelectorCache(b *testing.B) {
	DisableSelectorCache = true
	for i := 0; i < b.N; i++ {
		getQuery("/AAA/BBB/DDD/CCC/EEE/ancestor::*")
	}
}

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

func TestSelectorCache(t *testing.T) {
	SelectorCacheMaxEntries = 2
	for i := 1; i <= 3; i++ {
		getQuery(fmt.Sprintf("//a[position()=%d]", i))
	}
	getQuery("//a[position()=3]")

}

func TestLoadURLWithClient(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, htmlSample)
	}))
	defer ts.Close()
	cases := []struct {
		Name   string
		Client *http.Client
	}{
		{Name: "Client nil", Client: nil},
		{Name: "Client Custom", Client: &http.Client{Timeout: 10 * time.Second}},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			_, err := LoadURLWithClient(ts.URL, c.Client)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestLoadURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, htmlSample)
	}))
	defer ts.Close()

	_, err := LoadURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadURLWithGzipResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)
		defer gz.Close()
		fmt.Fprint(gz, htmlSample)
	}))
	defer ts.Close()

	_, err := LoadURL(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadDoc(t *testing.T) {
	tempHTMLdoc, err := ioutil.TempFile("", "sample_*.html")
	if err != nil {
		t.Fatal(err)
	}
	tempHTMLFilename := tempHTMLdoc.Name()
	defer func(tempHTMLdoc *os.File, filename string) {
		tempHTMLdoc.Close()
		os.Remove(filename)
	}(tempHTMLdoc, tempHTMLFilename)

	tempHTMLdoc.Write([]byte(htmlSample))

	if _, err := LoadDoc(tempHTMLFilename); err != nil {
		t.Fatal(err)
	}
}

func TestNavigator(t *testing.T) {
	top := FindOne(testDoc, "//html")
	nav := &NodeNavigator{curr: top, root: top, attr: -1}
	nav.MoveToChild() // HEAD
	nav.MoveToNext()
	if nav.NodeType() != xpath.TextNode {
		t.Fatalf("expectd node type is TextNode,but got %vs", nav.NodeType())
	}
	nav.MoveToNext() // <BODY>
	if nav.Value() != InnerText(FindOne(testDoc, "//body")) {
		t.Fatal("body not equal")
	}
	nav.MoveToPrevious() //
	nav.MoveToParent()   //<HTML>
	if nav.curr != top {
		t.Fatal("current node is not html node")
	}
	nav.MoveToNextAttribute()
	if nav.LocalName() != "lang" {
		t.Fatal("node not move to lang attribute")
	}

	nav.MoveToParent()
	nav.MoveToFirst() // <!DOCTYPE html>
	if nav.curr.Type != html.DoctypeNode {
		t.Fatalf("expected node type is DoctypeNode,but got %d", nav.curr.Type)
	}
}

func TestXPath(t *testing.T) {
	node := FindOne(testDoc, "//html")
	if SelectAttr(node, "lang") != "en-US" {
		t.Fatal("//html[@lang] != en-Us")
	}

	node = FindOne(testDoc, "//header")
	if strings.Index(InnerText(node), "Logo") > 0 {
		t.Fatal("InnerText() have comment node text")
	}
	if !strings.Contains(OutputHTML(node, true), "Logo") {
		t.Fatal("OutputHTML() shoud have comment node text")
	}
	link := FindOne(testDoc, "//a[1]/@href")
	if link == nil {
		t.Fatal("link is nil")
	}
	if v := InnerText(link); v != "/London" {
		t.Fatalf("expect value is /London, but got %s", v)
	}

}

func TestXPathCdUp(t *testing.T) {
	doc := loadHTML(`<html><b attr="1"></b></html>`)
	node := FindOne(doc, "//b/@attr/..")
	t.Logf("node = %#v", node)
	if node == nil || node.Data != "b" {
		t.Fatal("//b/@id/.. != <b></b>")
	}
}

func TestConcurrentQuery(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s := `<html><head></head><body><div>a</div></body>`
			doc := loadHTML(s)
			if n := FindOne(doc, `//div`); n == nil {
				t.Fatalf("should find one but got nil [%d]", i)
			}
		}(i)
	}
	wg.Wait()
}

func TestLoadURL_ErrorCases(t *testing.T) {
	tests := []struct {
		name            string
		contentEncoding string
		responseBody    []byte
		expectedError   string
	}{
		{
			name:            "ZlibReaderError",
			contentEncoding: "deflate",
			responseBody:    []byte("invalid zlib data"),
			expectedError:   "zlib",
		},
		{
			name:            "EmptyContentEncodingNilBody",
			contentEncoding: "",
			responseBody:    nil,
			expectedError:   "EOF",
		},
		{
			name:            "GzipReaderError",
			contentEncoding: "gzip",
			responseBody:    []byte("invalid gzip data"),
			expectedError:   "gzip",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a local HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Encoding", tt.contentEncoding)
				w.WriteHeader(http.StatusOK)
				w.Write(tt.responseBody)
			}))
			defer server.Close()

			node, err := LoadURL(server.URL)

			// Assert that the returned node is nil
			if node != nil {
				t.Errorf("Expected node to be nil, but got %v", node)
			}

			// Assert that an error is returned
			if err == nil {
				t.Error("Expected an error, but got nil")
			} else if !strings.Contains(err.Error(), tt.expectedError) {
				t.Errorf("Expected error to contain '%s', but got: %v", tt.expectedError, err)
			}
		})
	}
}

func loadHTML(str string) *html.Node {
	node, err := Parse(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	return node
}

func TestXPathDocumentStringValue(t *testing.T) {
	doc, err := Parse(strings.NewReader(`<!DOCTYPE html><html><head><title>A</title></head><body>B<!--excluded--><p>C</p></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		expr string
		want interface{}
	}{{`string(/)`, "ABC"}, {`string()`, "ABC"}, {`contains(/, 'BC')`, true}, {`string(//p)`, "C"}, {`string(//comment())`, "excluded"}} {
		t.Run(tc.expr, func(t *testing.T) {
			got := xpath.MustCompile(tc.expr).Evaluate(CreateXPathNavigator(doc))
			if got != tc.want {
				t.Errorf("got %#v want %#v", got, tc.want)
			}
		})
	}
}

// TestAttributeAxisCDUP verifies that the parent axis from an attribute
// returns the owning element, and that the navigator is no longer in
// attribute mode after the cdup.
func TestAttributeAxisCDUP(t *testing.T) {
	doc := loadHTML(`<html><div id="main" class="red" title="x"></div></html>`)
	nav := xpath.MustCompile(`//div/@id`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("attribute @id not found")
	}
	cur := nav.Current().(*NodeNavigator)
	if cur.NodeType() != xpath.AttributeNode {
		t.Fatalf("expected AttributeNode, got %v", cur.NodeType())
	}
	// Move to parent (cdup) from attribute
	if !cur.MoveToParent() {
		t.Fatal("MoveToParent from attribute should succeed")
	}
	if cur.NodeType() != xpath.ElementNode {
		t.Fatalf("after cdup expected ElementNode, got %v", cur.NodeType())
	}
	if cur.LocalName() != "div" {
		t.Fatalf("after cdup expected element div, got %q", cur.LocalName())
	}
}

// TestAttributeToNextAttribute covers iterating through all attributes
// and verifies that the navigator properly exits attribute mode once
// there are no more attributes.
func TestAttributeToNextAttribute(t *testing.T) {
	doc := loadHTML(`<html><a href="/x" rel="nofollow" target="_blank"></a></html>`)
	nav := xpath.MustCompile(`//a/@href`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("attribute @href not found")
	}
	cur := nav.Current().(*NodeNavigator)
	if cur.LocalName() != "href" || cur.Value() != "/x" {
		t.Fatalf("expected href=/x, got %s=%q", cur.LocalName(), cur.Value())
	}

	// Move to next attribute: should now be @rel
	if !cur.MoveToNextAttribute() {
		t.Fatal("expected next attribute @rel to exist")
	}
	if cur.LocalName() != "rel" || cur.Value() != "nofollow" {
		t.Fatalf("expected rel=nofollow, got %s=%q", cur.LocalName(), cur.Value())
	}

	// Move to next attribute: should now be @target
	if !cur.MoveToNextAttribute() {
		t.Fatal("expected next attribute @target to exist")
	}
	if cur.LocalName() != "target" || cur.Value() != "_blank" {
		t.Fatalf("expected target=_blank, got %s=%q", cur.LocalName(), cur.Value())
	}

	// Move past the last attribute: should return false and exit attribute mode
	if cur.MoveToNextAttribute() {
		t.Fatal("expected no more attributes")
	}
	// After returning false, navigator must no longer be in attribute mode
	if cur.NodeType() != xpath.ElementNode {
		t.Errorf("after exhausting attributes expected ElementNode, got %v", cur.NodeType())
	}
	if cur.LocalName() != "a" {
		t.Errorf("after exhausting attributes expected element a, got %q", cur.LocalName())
	}
}

// TestMoveToChildFromAttribute verifies that attempting to move to a child
// while positioned on an attribute returns false and clears the attribute
// cursor, so that the navigator is no longer in attribute mode.
func TestMoveToChildFromAttribute(t *testing.T) {
	doc := loadHTML(`<b attr="1"><i></i></b>`)
	iter := xpath.MustCompile("//b/@attr").Select(CreateXPathNavigator(doc))
	if !iter.MoveNext() {
		t.Fatal("attribute not found")
	}
	nav := iter.Current().(*NodeNavigator)

	if nav.MoveToChild() {
		t.Fatal("MoveToChild from attribute should return false")
	}
	if nav.attr != -1 {
		t.Fatalf("MoveToChild from attribute should reset attr to -1, got %d", nav.attr)
	}
	if nav.NodeType() != xpath.ElementNode {
		t.Fatalf("after reset, NodeType should be ElementNode, got %v", nav.NodeType())
	}
	if nav.LocalName() != "b" {
		t.Fatalf("after reset, LocalName should be 'b', got %q", nav.LocalName())
	}
}

// TestMoveToSiblingFromAttribute verifies that sibling moves while on an
// attribute return false and clear the attribute cursor.
func TestMoveToSiblingFromAttribute(t *testing.T) {
	doc := loadHTML(`<html><p class="x"></p><p class="y"></p></html>`)
	nav := xpath.MustCompile(`//p[1]/@class`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("attribute @class not found")
	}
	cur := nav.Current().(*NodeNavigator)

	// Next sibling
	if cur.MoveToNext() {
		t.Fatal("MoveToNext from attribute should return false")
	}
	if cur.NodeType() != xpath.ElementNode {
		t.Errorf("after MoveToNext from attribute expected ElementNode, got %v", cur.NodeType())
	}

	// Previous sibling
	if cur.MoveToPrevious() {
		t.Fatal("MoveToPrevious from attribute should return false")
	}
	if cur.NodeType() != xpath.ElementNode {
		t.Errorf("after MoveToPrevious from attribute expected ElementNode, got %v", cur.NodeType())
	}
}

// TestMoveToFirstFromAttribute verifies MoveToFirst from an attribute
// clears the attribute cursor and returns false.
func TestMoveToFirstFromAttribute(t *testing.T) {
	doc := loadHTML(`<html><b class="x" title="y"></b></html>`)
	iter := xpath.MustCompile("//b/@class").Select(CreateXPathNavigator(doc))
	if !iter.MoveNext() {
		t.Fatal("attribute not found")
	}
	nav := iter.Current().(*NodeNavigator)

	if nav.MoveToFirst() {
		t.Fatal("MoveToFirst from attribute should return false")
	}
	if nav.attr != -1 {
		t.Fatalf("MoveToFirst from attribute should reset attr to -1, got %d", nav.attr)
	}
	if nav.NodeType() != xpath.ElementNode {
		t.Fatalf("after reset, NodeType should be ElementNode, got %v", nav.NodeType())
	}
	if nav.LocalName() != "b" {
		t.Fatalf("after reset, LocalName should be 'b', got %q", nav.LocalName())
	}
}

// TestIdentityUniqueness verifies that Identity returns distinct values
// for different nodes and for different attributes, and that after
// leaving attribute mode it returns the element's identity.
func TestIdentityUniqueness(t *testing.T) {
	doc := loadHTML(`<html><p id="1" class="c"></p><p id="2"></p></html>`)
	nav := xpath.MustCompile(`//p`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("first p not found")
	}
	first := nav.Current().(*NodeNavigator)
	id1 := first.Identity()

	if !nav.MoveNext() {
		t.Fatal("second p not found")
	}
	second := nav.Current().(*NodeNavigator)
	id2 := second.Identity()

	if id1 == id2 {
		t.Errorf("expected distinct identities for different elements, got same %d", id1)
	}

	// Now move to an attribute of the first element and check identity differs
	first2 := xpath.MustCompile(`//p[1]/@id`).Select(CreateXPathNavigator(doc))
	if !first2.MoveNext() {
		t.Fatal("@id not found")
	}
	attrNav := first2.Current().(*NodeNavigator)
	attrID := attrNav.Identity()
	if attrID == id1 {
		t.Errorf("expected attribute identity to differ from element identity")
	}
}

// TestCopyPreservesAttributeState verifies that Copy duplicates the
// attribute cursor state so that the copy can continue iterating
// attributes independently.
func TestCopyPreservesAttributeState(t *testing.T) {
	doc := loadHTML(`<html><img src="a" alt="b" width="c"/></html>`)
	nav := xpath.MustCompile(`//img/@src`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("@src not found")
	}
	cur := nav.Current().(*NodeNavigator)

	cp := cur.Copy().(*NodeNavigator)
	if cp.LocalName() != "src" || cp.Value() != "a" {
		t.Fatalf("copy should preserve attribute, got %s=%q", cp.LocalName(), cp.Value())
	}

	// Advance the copy to the next attribute; original should be unaffected
	if !cp.MoveToNextAttribute() {
		t.Fatal("copy should be able to move to next attribute")
	}
	if cp.LocalName() != "alt" {
		t.Errorf("copy advanced to %q, expected alt", cp.LocalName())
	}
	if cur.LocalName() != "src" {
		t.Errorf("original should remain on src, got %q", cur.LocalName())
	}
}

// TestElementWithoutAttributes verifies that MoveToNextAttribute on an
// element with no attributes returns false and keeps the navigator in
// element mode (does not corrupt state).
func TestElementWithoutAttributes(t *testing.T) {
	doc := loadHTML(`<html><p></p></html>`)
	nav := xpath.MustCompile(`//p`).Select(CreateXPathNavigator(doc))
	if !nav.MoveNext() {
		t.Fatal("p not found")
	}
	cur := nav.Current().(*NodeNavigator)
	if cur.MoveToNextAttribute() {
		t.Fatal("element without attributes should have no attributes to move to")
	}
	if cur.NodeType() != xpath.ElementNode {
		t.Errorf("after failed MoveToNextAttribute expected ElementNode, got %v", cur.NodeType())
	}
}

// TestMoveToRestoresAttributeState verifies that MoveTo restores the full
// state (including the attribute cursor) of the source navigator, which
// is required for XPath engine backtracking.
func TestMoveToRestoresAttributeState(t *testing.T) {
	doc := loadHTML(`<html><a href="/1" title="t"></a></html>`)
	src := xpath.MustCompile(`//a/@href`).Select(CreateXPathNavigator(doc))
	if !src.MoveNext() {
		t.Fatal("@href not found")
	}
	srcNav := src.Current().(*NodeNavigator)

	dst := CreateXPathNavigator(doc)
	if !dst.MoveTo(srcNav) {
		t.Fatal("MoveTo should succeed for navigators of the same root")
	}
	if dst.LocalName() != "href" || dst.Value() != "/1" {
		t.Errorf("MoveTo should restore attribute state, got %s=%q", dst.LocalName(), dst.Value())
	}
	if dst.NodeType() != xpath.AttributeNode {
		t.Errorf("MoveTo should restore attribute mode, got %v", dst.NodeType())
	}

	// Moving to a navigator on a different root should fail
	other := loadHTML(`<html></html>`)
	otherNav := CreateXPathNavigator(other)
	if otherNav.MoveTo(srcNav) {
		t.Error("MoveTo across different roots should return false")
	}
}

// TestChainedAttributeAndChildAxis verifies the full chain:
// element -> attribute -> parent (cdup) -> child, ensuring the
// attribute cursor is properly cleared after cdup.
func TestChainedAttributeAndChildAxis(t *testing.T) {
	doc := loadHTML(`<html><ul class="list"><li>1</li><li>2</li></ul></html>`)
	node := FindOne(doc, `//ul/@class/..`)
	if node == nil {
		t.Fatal("cdup from @class should return the ul element")
	}
	if node.Data != "ul" {
		t.Errorf("expected ul, got %s", node.Data)
	}
	if SelectAttr(node, "class") != "list" {
		t.Errorf("expected class=list on cdup result, got %q", SelectAttr(node, "class"))
	}
}
