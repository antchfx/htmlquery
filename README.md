htmlquery
====
[![Build Status](https://travis-ci.org/antchfx/htmlquery.svg?branch=master)](https://travis-ci.org/antchfx/htmlquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/htmlquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/htmlquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/htmlquery?status.svg)](https://godoc.org/github.com/antchfx/htmlquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/htmlquery)](https://goreportcard.com/report/github.com/antchfx/htmlquery)

Overview
====

**htmlquery** is an HTML Parser and XPath package, supports extract data or evaluate from HTML documents using XPath expression.

[xmlquery](https://github.com/antchfx/xmlquery) is similar to this package, but supports XML document using XPath expression.

Installation
====

$ go get github.com/antchfx/htmlquery

Dependencies
====

- [html](https://golang.org/x/net/html)
- [xpath](https://github.com/antchfx/xpath)

Getting started
====

#### Load HTML document from URL

```go
doc, _ := htmlquery.LoadURL("http://example.com/")
fmt.Println(htmlquery.OutputHTML(doc, true))
```

#### Load HTML document from string

```go
s := `<html>....</html>`
doc, _ := htmlquery.Parse(strings.NewReader(s))
fmt.Println(htmlquery.OutputHTML(doc, true))
```

#### List all matched elements

```go
doc := loadTestHTML()
for _, n := range htmlquery.Find(doc, "//a") {
	fmt.Printf("%s \n", n.Data)
}
```

#### List all matched elements with specific attribute

```go
doc := loadTestHTML()
for _, n := range htmlquery.Find(doc, "//a/@href") {
	fmt.Printf("%s \n", htmlquery.SelectAttr(n, "href"))
}
```

#### Select first of all matched elements

```go
doc := loadTestHTML()
n := htmlquery.FindOne(doc, "//a")
fmt.Printf("%s \n", htmlquery.OutputHTML(n, true))
```

### Gets an element text or attribute value

```go
doc := loadTestHTML()
n := htmlquery.FindOne(doc, "//a")
fmt.Printf("%s \n", htmlquery.InnerText)
fmt.Println("%s \n", htmlquery.SelectAttr(n, "href"))
```

#### Evaluate element count

```go
doc := loadTestHTML()
expr, _ := xpath.Compile("count(//img)")
v := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64)
fmt.Printf("total count is %f", v)
```

Questions
===
If you have any questions, create an issue and welcome to contribute.