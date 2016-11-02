package main

import (
	"fmt"

	"github.com/antchfx/xquery/html"
)

func main() {
	u := "https://www.amazon.com/Microsoft-Surface-Pro-QG2-00001-4300U/dp/B00NVXE45U"
	doc, err := htmlquery.Load(u)
	if err != nil {
		panic(err)
	}
	var product struct {
		name     string
		price    string
		brand    string
		merchant string
	}
	product.name = htmlquery.InnerText(htmlquery.FindOne(doc, "//span[@id='productTitle']"))
	product.brand = htmlquery.InnerText(htmlquery.FindOne(doc, "//a[@id='brand']"))
	product.price = htmlquery.InnerText(htmlquery.FindOne(doc, "//span[@id='priceblock_ourprice']"))
	product.merchant = htmlquery.InnerText(htmlquery.FindOne(doc, "//div[@id='merchant-info']/a"))

	fmt.Println(product.name)
	fmt.Println(product.brand)
	fmt.Println(product.price)
	fmt.Println(product.merchant)
}
