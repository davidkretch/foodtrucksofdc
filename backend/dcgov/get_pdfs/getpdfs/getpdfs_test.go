package getpdfs

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestGet(t *testing.T) {
	resp, err := GetURL("https://httpbin.org/get")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("Get returned status code != 200")
	}
}

func TestGetLinks(t *testing.T) {
	doc := strings.NewReader("<html><body><a href=\"url0\">text0</a><a href=\"url1\">text1</a></body></html>")
	links := GetLinks(doc)

	if len(links) != 2 {
		t.Fatalf("GetLinks returned incorrect number of links")
	}

	if links[0].URL != "url0" {
		t.Fatalf("GetLinks returned wrong URL")
	}
	if links[0].Text != "text0" {
		t.Fatalf("GetLinks returned wrong text")
	}

	if links[1].URL != "url1" {
		t.Fatalf("GetLinks returned wrong URL")
	}
	if links[1].Text != "text1" {
		t.Fatalf("GetLinks returned wrong text")
	}
}

func TestGetAttribute(t *testing.T) {
	token := html.Token{
		Attr: []html.Attribute{
			html.Attribute{Namespace: "", Key: "foo", Val: "bar"},
			html.Attribute{Namespace: "", Key: "baz", Val: "qux"},
		},
	}

	attr, err := GetAttribute(token, "foo")
	if err != nil {
		t.Fatalf("GetAttribute threw an unexpected error: %v", err)
	}
	if attr != "bar" {
		t.Fatalf("GetAttribute returned an incorrect value: %s", attr)
	}

	attr, err = GetAttribute(token, "baz")
	if err != nil {
		t.Fatalf("GetAttribute threw an unexpected error: %v", err)
	}
	if attr != "qux" {
		t.Fatalf("GetAttribute returned an incorrect value: %s", attr)
	}

	attr, err = GetAttribute(token, "invalid")
	if err == nil {
		t.Fatalf("GetAttribute failed to throw an error for an invalid key")
	}
}

func TestFilter(t *testing.T) {
	links := []Link{
		Link{Text: "foo0", URL: "bar0"},
		Link{Text: "foo1", URL: "bar1"},
		Link{Text: "baz", URL: "qux"},
	}

	r := Filter(links, func(Link) bool { return true })
	if len(r) != len(links) {
		t.Fatal("Filter(_, true) returned other than all elements")
	}

	r = Filter(links, func(Link) bool { return false })
	if len(r) != 0 {
		t.Fatal("Filter(_, false) returned other than zero elements")
	}

	r = Filter(links, func(l Link) bool {
		return strings.HasPrefix(l.Text, "foo")
	})
	if len(r) != 2 {
		t.Fatal("Filter(_, HasPrefix('foo')) returned other than 2 elements")
	}

	r = Filter(links, func(l Link) bool { return l.Text == "baz" })
	if len(r) != 1 {
		t.Fatal("Filter(_, 'baz') returned other than 1 elements")
	}
}
