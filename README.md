XQuery
====
[![Build Status](https://travis-ci.org/antchfx/xquery.svg?branch=master)](https://travis-ci.org/antchfx/xquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/xquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/xquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/xquery?status.svg)](https://godoc.org/github.com/antchfx/xquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/xquery)](https://goreportcard.com/report/github.com/antchfx/xquery)

A golang package, extract data from HTML and XML documents using XPath expression.

Supported most of XPath feature(syntax), you can found at [XPath](https://github.com/antchfx/xpath).

Installing
====

> go get github.com/antchfx/xquery

HTML Query
===

Extract data from HTML document using XPath.

[![GoDoc](https://godoc.org/github.com/antchfx/xquery/html?status.svg)](https://godoc.org/github.com/antchfx/xquery/html)

```go
package main

import (
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
	root, err := htmlquery.Parse(strings.NewReader(s))
	if err != nil {
		panic(err)
	}
	node := htmlquery.FindOne(root, "//title")
	fmt.Println(htmlquery.InnerText(node))	
}
```

XML Query
===
Extract data from XML document using XPath.

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
