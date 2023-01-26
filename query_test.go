package htmlquery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

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

func TestSelectorCache(t *testing.T) {
	SelectorCacheMaxEntries = 2
	for i := 1; i <= 3; i++ {
		getQuery(fmt.Sprintf("//a[position()=%d]", i))
	}
	getQuery("//a[position()=3]")

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
	wg.Done()
}

func loadHTML(str string) *html.Node {
	node, err := Parse(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	return node
}
