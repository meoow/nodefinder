#Node Finder

This library finds all nodes in html file given specific path.

##Usage

####path example:

```
//this path:
body/div.good[id="a b",code="go"]/#gogogo/p.download[os=linux,arch=x86_64]/a[href="http://meow.com/",target=_blank]
```
```html
<!-- will be parsed into -->
<body><div class="good" id="a b" code="go">< id="gogogo"><p class="download" os="linux" arch="x86_64"><a href="http://meow.com/" target="_blank">
```

####code Example:
```go
package main

import (
	"net/http"
	"github.com/meoow/nodefinder"
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
	tags := nodefinder.TagParser(path)
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

```