#Node Finder

This library finds all nodes in html file given specific path.

##Usage

####syntax:
`body/div.class1/div#id1//table[a=b sometag='c1 c2 c3']/tbody.class2//td.class2:3`  
1. **body** is a single tag. If preceded by "/", it forces to match the root node.  
2. **/** matches direct children.  
3. **div.class1** is a div tag with class="class1".   
4. **div#id1** div tag with id="id1".  
5. **//** matches children in any depth. Leading // in path will be ignored.  
6. **table[a=b c=good sometag='c1 c2 c3' ]** attributes key-value pairs are listed in square brackets, if multiple values are needed, put them in single/double quotes, the order is insignificant.  
7. **td:3** matches the 3rd occurence of tag td with class="class2".  

####path example:

```
//this path:
body/div.good[id="a b" code=go]/#gogogo/p.download[os=linux arch=x86_64]/a[href="http://meow.com/",target=_blank]
```
```html
<!-- will match this -->
<body><div class="good" id="a b" code="go"><div id="gogogo"><p class="download" os="linux" arch="x86_64"><a href="http://meow.com/" target="_blank">
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

```
