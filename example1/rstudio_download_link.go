package main

import (
	"bytes"
	"fmt"
	"github.com/meoow/nodefinder"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	link := "http://www.rstudio.com/products/rstudio/download/"
	path := "table.downloads/tbody/tr/td/a"
	md5path := "table.downloads/tbody/tr/td/a/../../td/code"

	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	contentsReader := bytes.NewReader(contents)

	fmt.Fprintln(os.Stderr, "Start Parsing...")

	tags := nodefinder.NewPath(path)
	md5tags := nodefinder.NewPath(md5path)

	md5nodes, _ := nodefinder.Find(md5tags, contentsReader)

	contentsReader.Seek(0, 0)

	nodes, _ := nodefinder.Find(tags, contentsReader)

	printLinkOnly := false
	if len(nodes) != len(md5nodes) {
		printLinkOnly = true
	}

	for i, n := range nodes {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val)
			}
		}
		if !printLinkOnly {
			fmt.Println(md5nodes[i].FirstChild.Data)
		}
	}
}
