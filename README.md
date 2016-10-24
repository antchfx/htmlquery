XQuery
====
XQuery is a package to extract data from HTML and XML using XPath selectors.

* [HTML](https://godoc.org/golang.org/x/net/html) : The HTML parser package,from golang official.

* [XML](https://github.com/antchfx/xml) : The XML parser package.

* [GXPATH](https://github.com/antchfx/gxpath) : The XPath package for Go.

Installing
====

> go get -u github.com/antchfx/xquery

#### HTML Usage

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
	title := htmlquery.FindOne(root, "//title")
	fmt.Println(fmt.Sprintf("document title: %s", htmlquery.InnerText(title)))
}
```

#### XML Usage

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
	title := xmlquery.FindOne(root, "//title")
	fmt.Println(fmt.Sprintf("document title: %s", title.InnerText()))
}
```