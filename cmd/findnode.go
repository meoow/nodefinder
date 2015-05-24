package main

import (
	"flag"
	"fmt"
	"github.com/meoow/nodefinder"
	"golang.org/x/net/html"
	"os"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "%s -i HTML PATH [PATH ...]\n", os.Args[0])
	flag.PrintDefaults()
}

var input string
var inner bool

func init() {
	flag.Usage = usage
	flag.StringVar(&input, "i", "", "input html file")
	flag.BoolVar(&inner, "I", false, "print inner contents of matched tags only")
}

func main() {
	flag.Parse()

	if input == "" {
		flag.Usage()
		os.Exit(1)
	}

	if fi, err := os.Stat(input); err != nil || fi.IsDir() {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Fprintf(os.Stderr, "%s is not a file", fi.Name())
		}
		flag.Usage()
		os.Exit(1)
	}

	var xpath string
	if flag.NArg() > 0 {
		xpath = strings.Join(flag.Args(), ":::")
	}

	fh, err := os.Open(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tags := nodefinder.NewPath(xpath)
	nodes, err := nodefinder.Find(tags, fh)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	for _, i := range nodes {
		if inner {
			for j := i.FirstChild; j != nil; j = j.NextSibling {
				html.Render(os.Stdout, j)
			}
		} else {
			html.Render(os.Stdout, i)
		}
	}
}
