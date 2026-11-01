package htmlquery

import (
	"github.com/antchfx/xpath"
	"strings"
	"testing"
)

func TestXPathDocumentStringValue(t *testing.T) {
	doc, err := Parse(strings.NewReader(`<!DOCTYPE html><html><head><title>A</title></head><body>B<!--excluded--><p>C</p></body></html>`))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		expr string
		want interface{}
	}{{`string(/)`, "ABC"}, {`string()`, "ABC"}, {`contains(/, 'BC')`, true}, {`string(//p)`, "C"}, {`string(//comment())`, "excluded"}} {
		t.Run(tc.expr, func(t *testing.T) {
			got := xpath.MustCompile(tc.expr).Evaluate(CreateXPathNavigator(doc))
			if got != tc.want {
				t.Errorf("got %#v want %#v", got, tc.want)
			}
		})
	}
}
