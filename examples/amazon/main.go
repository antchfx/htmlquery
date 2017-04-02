package main

import (
	"fmt"
	"strings"

	"github.com/antchfx/xquery/html"
)

func main() {
	u := "https://www.amazon.com/Microsoft-Surface-Pro-QG2-00001-4300U/dp/B00NVXE45U"
	doc, err := htmlquery.LoadURL(u)
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

	fmt.Printf("%s \n", u)
	fmt.Printf("Name: %s\n", strings.TrimSpace(product.name))
	fmt.Printf("Brand: %s\n", product.brand)
	fmt.Printf("Price: %s\n", product.price)
	fmt.Printf("Merchant: %s\n", product.merchant)
}
