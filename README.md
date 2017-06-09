XQuery
====
[![Build Status](https://travis-ci.org/antchfx/xquery.svg?branch=master)](https://travis-ci.org/antchfx/xquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/xquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/xquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/xquery?status.svg)](https://godoc.org/github.com/antchfx/xquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/xquery)](https://goreportcard.com/report/github.com/antchfx/xquery)

A golang package, lets you extract data from HTML/XML documents using XPath expression.

List of supported XPath functions you can found at [XPath Package](https://github.com/antchfx/xpath).

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
	f, _ := os.Open(`./examples/test.html`)
	doc, _ := htmlquery.Parse(f)
	expr := xpath.MustCompile("count(//div[@class='article'])")
	fmt.Printf("%f \n", expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64))

	expr = xpath.MustCompile("//a/@href")
	iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
	for iter.MoveNext() {
		fmt.Printf("%s \n", iter.Current().Value()) // output href
	}

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
	f, _ := os.Open(`./examples/test.xml`)
	doc, _ := xmlquery.Parse(f)
	// sum all book's price via Evaluate()
	expr, err := xpath.Compile("sum(//book/price)")
	if err != nil {
		panic(err)
	}
	fmt.Printf("total price: %f\n", expr.Evaluate(xmlquery.CreateXPathNavigator(doc)).(float64))

	for _, n := range xmlquery.Find(doc, "//book") {
		fmt.Printf("%s : %s \n", n.SelectAttr("id"), xmlquery.FindOne(n, "title").InnerText())
	}
	
	n := xmlquery.FindOne(doc, "//book[@id='bk104']")
	fmt.Printf("%s \n", n.OutputXML(true))
}
```
