package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func readHTML(filename string) (htmltext, version string) {
	fh, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fh.Close()
	root, err := html.Parse(fh)
	if err != nil {
		log.Fatal(err)
	}
	version = walkHTML(root)
	var sb strings.Builder
	if err = html.Render(&sb, root); err != nil {
		log.Fatal(err)
	}
	return sb.String(), version
}

func walkHTML(node *html.Node) (version string) {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if v := walkHTML(child); v != "" {
			version = v
		}
	}
	if node.Type == html.ElementNode && node.DataAtom == atom.Div {
		for i, a := range node.Attr {
			if a.Key == "class" && a.Val == "version" {
				if node.FirstChild == nil || node.FirstChild.Type != html.TextNode || node.FirstChild.NextSibling != nil {
					log.Fatal("unexpected format for div version")
				}
				version = node.FirstChild.Data
			}
			if a.Key == "data-include-html" {
				filename := filepath.Join("resources", "html", a.Val+".html")
				fh, err := os.Open(filename)
				if err != nil {
					log.Fatal(err)
				}
				defer fh.Close()
				root, err := html.Parse(fh)
				if err != nil {
					log.Fatal(err)
				}
				node.FirstChild, node.LastChild = root, root
				node.Attr[i].Key = "data-included-html"
				// Note that we do this after visiting children,
				// so we do not traverse into the include
				// looking for nested includes; we also don't
				// see version numbers in includes.
			}
		}
	}
	return version
}
