package main

import (
	"net/http"
	"sync"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	urlRoot = "http://ruliweb.com"
)

func parseMainNodes(n *html.Node) bool {
	if n.DataAtom == atom.A && n.Parent != nil {
		return scrape.Attr(n.Parent, "class") == "row"
	}
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

var wg sync.WaitGroup

func main() {
	resp, err := http.Get(urlRoot)
	errCheck(err)

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	errCheck(err)

	urlList := scrape.FindAll(root, parseMainNodes)
}
