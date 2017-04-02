package main

import (
	"fmt"

	"github.com/antchfx/xquery/html"
	"golang.org/x/net/html"
)

func main() {
	q := "golang"
	u := "http://www.bing.com/search?q=" + q
	doc, err := htmlquery.LoadURL(u)
	if err != nil {
		panic(err)
	}
	type entry struct {
		id    int
		title string
		url   string
		desc  string
	}

	var entries []entry
	htmlquery.FindEach(doc, "//ol[@id='b_results']/li[@class='b_algo']", func(i int, node *html.Node) {
		item := entry{}
		item.id = i
		h2 := htmlquery.FindOne(node, "//h2")
		item.title = htmlquery.InnerText(h2)
		item.url = htmlquery.SelectAttr(htmlquery.FindOne(h2, "a"), "href")
		if n := htmlquery.FindOne(node, "//div[@class='b_caption']/p"); n != nil {
			item.desc = htmlquery.InnerText(n)
		}
		entries = append(entries, item)
	})
	count := htmlquery.InnerText(htmlquery.FindOne(doc, "//span[@class='sb_count']"))
	fmt.Println(fmt.Sprintf("%s by %s", count, q))
	for _, item := range entries {
		fmt.Println(fmt.Sprintf("%d title: %s", item.id, item.title))
		fmt.Println(fmt.Sprintf("url: %s", item.url))
		fmt.Println(fmt.Sprintf("desc: %s", item.desc))
		fmt.Println("=====================")
	}
}
