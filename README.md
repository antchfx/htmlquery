XQuery
====
XQuery is a package to extract data from HTML and XML using XPath selectors.

* [HTML](https://godoc.org/golang.org/x/net/html) : The HTML parser package,from golang official.

* [XML](https://github.com/antchfx/xml) : The XML parser package.

* [GXPATH](https://github.com/antchfx/gxpath) : The XPath package for Go.

Installing
====

> go get -u github.com/antchfx/xquery

#### HTML Query

Methods: 
* Find(*html.Node, string) []*html.Node
* FindOne(*html.Node, string) *html.Node
* FindEach(*html.Node, string, func(int, *html.Node))
* Load(string) *html.Node

```go
package main

import (
    "golang.org/x/net/html"
    "github.com/antchfx/xquery/html"	
)

func main() {	
	root, err := html.Parse(strings.NewReader(html_string))
	if err != nil {
		panic(err)
	}
	node := htmlquery.FindOne(root, "//title")
	fmt.Println(htmlquery.InnerText(node))
}
```

#### XML Query

Methods: 
* Find(*xml.Node, string) []*xml.Node
* FindOne(*xml.Node, string) *xml.Node
* FindEach(*xml.Node, string, func(int, *xml.Node))

```go
package main

import (
	"github.com/antchfx/xml"
	"github.com/antchfx/xquery/xml"
)

func main() {
	root, err := xml.Parse(strings.NewReader(xml_string))
	if err != nil {
		panic(err)
	}
	node := xmlquery.FindOne(root, "//title")
	fmt.Println(node.InnerText())
}
```