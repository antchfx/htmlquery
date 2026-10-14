package htmlquery

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"io"
	"net/http"
	"testing"
)

type trackedResponseBody struct {
	io.Reader
	closes int
}

func (b *trackedResponseBody) Close() error {
	b.closes++
	return nil
}

type responseTransport func(*http.Request) (*http.Response, error)

func (f responseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type failedResponseReader struct{}

func (failedResponseReader) Read([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestLoadURLWithClientClosesResponseBody(t *testing.T) {
	const document = "<html><head><title>example</title></head><body>page</body></html>"
	var gz, def bytes.Buffer
	gw := gzip.NewWriter(&gz)
	if _, err := io.WriteString(gw, document); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	dw := zlib.NewWriter(&def)
	if _, err := io.WriteString(dw, document); err != nil {
		t.Fatal(err)
	}
	if err := dw.Close(); err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name, encoding string
		reader         io.Reader
		wantErr        bool
	}{
		{"plain", "", bytes.NewBufferString(document), false},
		{"gzip", "gzip", bytes.NewReader(gz.Bytes()), false},
		{"deflate", "deflate", bytes.NewReader(def.Bytes()), false},
		{"invalid gzip", "gzip", bytes.NewBufferString("invalid"), true},
		{"invalid deflate", "deflate", bytes.NewBufferString("invalid"), true},
		{"unsupported encoding", "br", bytes.NewBufferString(document), true},
		{"read failure", "", failedResponseReader{}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := &trackedResponseBody{Reader: tc.reader}
			client := &http.Client{Transport: responseTransport(func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header: http.Header{
						"Content-Encoding": []string{tc.encoding},
						"Content-Type":     []string{"text/html; charset=utf-8"},
					},
					Body:    body,
					Request: req,
				}, nil
			})}

			node, err := LoadURLWithClient("https://example.invalid/page", client)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, want error %t", err, tc.wantErr)
			}
			if err == nil {
				title := FindOne(node, "//title")
				if title == nil || InnerText(title) != "example" {
					t.Fatal("parsed document changed")
				}
			}
			if body.closes != 1 {
				t.Errorf("response body Close calls = %d, want 1", body.closes)
			}
		})
	}
}
