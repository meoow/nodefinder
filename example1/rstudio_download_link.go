package main

import (
	"github.com/meoow/nodefinder"
	"net/http"
)
import "log"
import "fmt"

func main() {
	link := "http://www.rstudio.com/products/rstudio/download/"
	path := "table.downloads/tbody/tr/td/a"
	resp, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("Start Parsing...")
	tags := nodefinder.NewPath(path)
	nodes, err := nodefinder.Find(tags, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	for _, n := range nodes {
		for _, a := range n.Attr {
			if a.Key == "href" {
				fmt.Println(a.Val)
			}
		}
	}
}
