XQuery
====
XQuery is a golang package to extract data from HTML and XML using XPath selectors.

Supported most of XPath feature(syntax), you can found at [XPath](https://github.com/antchfx/xpath).

Installing
====

> go get -u github.com/antchfx/xquery

HTML Query
===

Lets extract data from HTML document using XPath.

[![GoDoc](https://godoc.org/github.com/antchfx/xquery/html?status.svg)](https://godoc.org/github.com/antchfx/xquery/html)

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
Lets extract data from XML document using XPath.

[![GoDoc](https://godoc.org/github.com/antchfx/xquery/xml?status.svg)](https://godoc.org/github.com/antchfx/xquery/xml)

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