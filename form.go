package main

import (
	"net/url"

	"github.com/andybalholm/cascadia"

	"golang.org/x/net/html"
)

var (
	selSubmit = cascadia.MustCompile("form.submitForm")
	selEnter  = cascadia.MustCompile("form#enterForm")
	selInput  = cascadia.MustCompile("input[name]")
)

func form(nodes []*html.Node) url.Values {
	f := url.Values{}
	for _, i := range nodes {
		var k, v string
		for _, a := range i.Attr {
			if a.Key == "name" {
				k = a.Val
			} else if a.Key == "value" {
				v = a.Val
			}
		}
		f.Add(k, v)
	}
	return f
}

func formAction(node *html.Node) string {
	for _, a := range node.Attr {
		if a.Key == "action" {
			return a.Val
		}
	}
	return ""
}
