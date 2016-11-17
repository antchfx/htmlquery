XQuery
====
XQuery is a package to extract data from HTML and XML using XPath selectors.

Installing
====

> go get -u github.com/antchfx/xquery

HTML Query
===
This package use golang official package to parse html document: [html](https://godoc.org/golang.org/x/net/html).

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
	html_string:=`<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<body>
<h1>This is a Heading</h1>
<p>This is a paragraph.</p>
</body>
</html>`
	root, err := html.Parse(strings.NewReader(html_string))
	if err != nil {
		panic(err)
	}
	node := htmlquery.FindOne(root, "//title")
	fmt.Println(htmlquery.OutputHTML(node)) // output html text with tags
	fmt.Println(htmlquery.InnerText(node))	
}
```

XML Query
===

Methods: 
* Find(*Node, string) []*Node
* FindOne(*Node, string) *Node
* FindEach(*Node, string, func(int, *Node))

```go
package main

import (
	"github.com/antchfx/xquery/xml"
)

func main() {
	xml_string:=`<?xml version="1.0" encoding="UTF-8"?>
<bookstore>
<book category="cooking">
  <title lang="en">Everyday Italian</title>
  <author>Giada De Laurentiis</author>
  <year>2005</year>
  <price>30.00</price>
</book>
......
</bookstore>`
	root, err := xmlquery.ParseXML(strings.NewReader(xml_string))
	if err != nil {
		panic(err)
	}
	node := xmlquery.FindOne(root, "//book[@category='cooking']")
	fmt.Println(node.OutputXML()) // output xml text with tags
	fmt.Println(node.InnerText()) // output text without tags
}
```