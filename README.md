XQuery
====
XQuery is a golang package to extract data from HTML and XML using XPath selectors.

Installing
====

> go get -u github.com/antchfx/xquery

HTML Query
===
This package use golang [html](https://godoc.org/golang.org/x/net/html) package to parse a HTML document.

|Method                    |Descript|
|--------------------------|----------------|
|LoadURL(url string) `*html.Node` |Loads the HTML document from the specified URL|
|Find(`*html.Node`, expr string) []`*html.Node`|Searches all the html.Node that matches the specified XPath expression expr|
|FindOne(`*html.Node`, expr string) `*html.Node`|Searches the html.Node and returns a first matched node|
|FindEach(`*html.Node`, expr string,cb func(int, `*html.Node`))|Searches all the matched html.Node and to pass its a callback function cb|
|OutputHTML(`*html.Node`) string|Returns html format output of this html.Node|
|InnerText(`*html.Node`) string|Returns text without html tag of this html.Node|
|SelectAttr(`*html.Node`, name string) string|Returns the attribute value with the specified attribute name|

```go
package main

import (
    "golang.org/x/net/html"
    "github.com/antchfx/xquery/html"	
)

func main() {
	s:=`<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<body>
<h1>This is a Heading</h1>
<p>This is a paragraph.</p>
</body>
</html>`
	root, err := html.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	node := htmlquery.FindOne(root, "//title")
	fmt.Println(htmlquery.InnerText(node))	
}
```

XML Query
===
This package is similar to the HTML query package, its implementation load XML document and parsing.

|Method                    |Descript|
|--------------------------|----------------|
|LoadURL(url string) (`*xmlquery.Node`, error) |Loads the XML document from the specified URL|
|Parse(io.Reader) (`*xmlquery.Node`, error)|Parses the specified io.Reader to the XML document.|
|Find(`*xmlquery.Node`, expr string) []`*xmlquery.Node`|Searches all the xmlquery.Node that matches the specified XPath expression expr|
|FindOne(`*xmlquery.Node`, expr string) `*xmlquery.Node`|Searches the xmlquery.Node and returns a first matched node|
|FindEach(`*xmlquery.Node`, expr string,cb func(int, `*xmlquery.Node`))|Searches all the matched Node and to pass its a callback function cb|
|Node.SelectElements(name string)[]`*xmlquery.Node`|Finds child elements with the specified element name|
|Node.SelectElement(name string)`*xmlquery.Node`|Finds child elements with the specified element name|
|Node.SelectAttr(name string)string|Returns the attribute value with the specified attribute name|
|Node.OutputHTML() string|Returns html format output of this node|
|Node.InnerText() string|Returns text without xml element tag of this Node|

```go
package main

import (
	"github.com/antchfx/xquery/xml"
)

func main() {
	s:=`<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
<book category="cooking">
  <title lang="en">Everyday Italian</title>
  <author>Giada De Laurentiis</author>
  <year>2005</year>
  <price>30.00</price>
</book>
......
</bookstore>`
	root, err := xmlquery.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	node := xmlquery.FindOne(root, "//book[@category='cooking']")
	fmt.Println(node.InnerText())
}
```