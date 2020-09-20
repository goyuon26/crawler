package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
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
	return false
}

func errCheck(err error) {
	if err != nil {
		panic(err)
	}
}

var wg sync.WaitGroup

func scrapContents(url string, fn string) {
	defer wg.Done()

	resp, err := http.Get(url)
	errCheck(err)

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	errCheck(err)

	matchNode := func(n *html.Node) bool {
		return n.DataAtom == atom.A && scrape.Attr(n, "class") == "deco"
	}

	file, err := os.OpenFile("c:/scrape/"+fn+".txt", os.O_CREATE|os.O_RDWR, os.FileMode(0777))
	errCheck(err)

	defer file.Close()

	w := bufio.NewWriter(file)

	for _, g := range scrape.FindAll(root, matchNode) {
		fmt.Println("result : ", scrape.Text(g))
		w.WriteString(scrape.Text(g) + "\r\n")
	}
	w.Flush()

}

func main() {
	resp, err := http.Get(urlRoot)
	errCheck(err)

	defer resp.Body.Close()

	root, err := html.Parse(resp.Body)
	errCheck(err)

	urlList := scrape.FindAll(root, parseMainNodes)

	for _, link := range urlList {

		//fmt.Println("Check Main Link : ", link, idx)
		//fmt.Println("Target URL : ", scrape.Attr(link, "href"))
		fileName := strings.Replace(scrape.Attr(link, "href"), "https://bbs.ruliweb.com/family/", "", 1)
		wg.Add(1)

		go scrapContents(scrape.Attr(link, "href"), fileName)
	}

	wg.Wait()
}
