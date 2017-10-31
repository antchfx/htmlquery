XQuery
====
[![Build Status](https://travis-ci.org/antchfx/xquery.svg?branch=master)](https://travis-ci.org/antchfx/xquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/xquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/xquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/xquery?status.svg)](https://godoc.org/github.com/antchfx/xquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/xquery)](https://goreportcard.com/report/github.com/antchfx/xquery)

A golang package, lets you extract data from HTML/XML documents using XPath expression.

List of supported XPath functions you can found here [XPath Package](https://github.com/antchfx/xpath).

Install
====

> go get github.com/antchfx/xquery


HTML Query [![GoDoc](https://godoc.org/github.com/antchfx/xquery/html?status.svg)](https://godoc.org/github.com/antchfx/xquery/html)
===

Extract data from HTML document.

```go
package main

import (
	"github.com/antchfx/xpath"
	"github.com/antchfx/xquery/html"
)

func main() {
	// Load HTML file.
	f, err := os.Open(`./examples/test.html`)
	if err != nil {
		panic(err)
	}
	// Parse HTML document.
	doc, err := htmlquery.Parse(f)
	if err != nil{
		panic(err)
	}

	// Option 1: using xpath's expr to matches nodes.
	expr := xpath.MustCompile("count(//div[@class='article'])")
	fmt.Printf("%f \n", expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64))

	expr = xpath.MustCompile("//a/@href")
	iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
	for iter.MoveNext() {
		fmt.Printf("%s \n", iter.Current().Value()) // output href
	}

	// Option 2: using build-in functions Find() to matches nodes.
	for _, n := range htmlquery.Find(doc, "//a/@href") {
		fmt.Printf("%s \n", htmlquery.SelectAttr(n, "href")) // output href
	}
}
```

XML Query [![GoDoc](https://godoc.org/github.com/antchfx/xquery/xml?status.svg)](https://godoc.org/github.com/antchfx/xquery/xml)
===
Extract data from XML document.

```go
package main

import (
	"github.com/antchfx/xpath"
	"github.com/antchfx/xquery/xml"
)

func main() {
	// Load XML document from file.
	f, err := os.Open(`./examples/test.xml`)
	if err != nil {
		panic(err)
	}
	// Parse XML document.
	doc, err := xmlquery.Parse(f)
	if err != nil{
		panic(err)
	}

	// Option 1: using xpath's expr to matches nodes.

	// sum all book's price via Evaluate()
	expr, err := xpath.Compile("sum(//book/price)")
	if err != nil {
		panic(err)
	}
	fmt.Printf("total price: %f\n", expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64))

	for _, n := range xmlquery.Find(doc, "//book") {
		fmt.Printf("%s : %s \n", n.SelectAttr("id"), xmlquery.FindOne(n, "title").InnerText())
	}
	
	// Option 2: using build-in functions FindOne() to matches node.
	n := xmlquery.FindOne(doc, "//book[@id='bk104']")
	fmt.Printf("%s \n", n.OutputXML(true))
}
```
