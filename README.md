htmlquery
====
[![Build Status](https://travis-ci.org/antchfx/htmlquery.svg?branch=master)](https://travis-ci.org/antchfx/htmlquery)
[![Coverage Status](https://coveralls.io/repos/github/antchfx/htmlquery/badge.svg?branch=master)](https://coveralls.io/github/antchfx/htmlquery?branch=master)
[![GoDoc](https://godoc.org/github.com/antchfx/htmlquery?status.svg)](https://godoc.org/github.com/antchfx/htmlquery)
[![Go Report Card](https://goreportcard.com/badge/github.com/antchfx/htmlquery)](https://goreportcard.com/report/github.com/antchfx/htmlquery)

htmlquery, lets you extract data from HTML documents using XPath expression, depend on 
[html](https://golang.org/x/net/html) and [xpath](https://github.com/antchfx/xpath) packages.

Installation
====

    $ go get github.com/antchfx/htmlquery

Examples
====

```go
func main() {
	// Load HTML file.
	f, err := os.Open(`./examples/test.html`)
	if err != nil {
		panic(err)
	}
	// Parse HTML document
	doc, err := htmlquery.Parse(f)
	if err != nil{
		panic(err)
	}

	// List all matches nodes with the name `a`.
	for _, n := range htmlquery.Find(doc, "//a") {
		fmt.Printf("%s \n", n.Data)
	}
}
```

### Evaluate an XPath expression

Using `Evaluate()` to evaluates XPath expressions.

```go
expr := xpath.MustCompile("count(//div[@class='article'])")
fmt.Printf("%f \n", expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(float64))

expr = xpath.MustCompile("//a/@href")
iter := expr.Evaluate(htmlquery.CreateXPathNavigator(doc)).(*xpath.NodeIterator)
for iter.MoveNext() {
	fmt.Printf("%s \n", iter.Current().Value())
}
```