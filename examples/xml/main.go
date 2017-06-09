package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/antchfx/xquery/xml"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := `<?xml version="1.0"?>
        <note>
		<to>Tove</to>
		<from>Jani</from>
		<heading>Reminder</heading>
		<body>Don't forget me this weekend!</body>
		</note>`
		w.Header().Set("Content-Type", "text/xml")
		w.Write([]byte(s))
	}))
	defer server.Close()

	root, err := xmlquery.LoadURL(server.URL)
	if err != nil {
		panic(err)
	}
	var n *xmlquery.Node
	n = xmlquery.FindOne(root, "//from")
	fmt.Println(fmt.Sprintf("from: %s", n.InnerText()))
	n = xmlquery.FindOne(root, "//to")
	fmt.Println(fmt.Sprintf("to: %s", n.InnerText()))
}
