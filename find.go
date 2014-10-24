package nodefinder

import "code.google.com/p/go.net/html"
import "io"
import "strings"

func Find(tags []*Tag, hf io.Reader) ([]*html.Node, error) {
	if len(tags) == 0 {
		return []*html.Node{}, nil
	}

	n, err := html.Parse(hf)
	if err != nil {
		return []*html.Node{}, err
	}

	roots := make([]*html.Node, 0, 1)
	result := make([]*html.Node, 0, 1)

	find1(tags[0], n, &roots)

	if len(tags) == 1 {
		for _, r := range roots {
			result = append(result, r)
		}
	} else {
		for _, f := range roots {
			find2(tags[1:len(tags)], f, &result)
		}
	}
	return result, nil
}

func find1(tag *Tag, p *html.Node, found *[]*html.Node) {
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if c.Data == tag.Tag {
			if compare(tag, c) {
				*found = append(*found, c)
			}
		}
		find1(tag, c, found)
	}
}

func find2(tags []*Tag, p *html.Node, result *[]*html.Node) {
	for _, cn := range childs(p) {
		if cn.Type != html.ElementNode {
			continue
		}
		if tags[0].Tag == cn.Data && compare(tags[0], cn) {
			if len(tags) == 1 {
				*result = append(*result, cn)
			} else {
				find2(tags[1:len(tags)], cn, result)
			}
		}
	}
}

func childs(n *html.Node) []*html.Node {
	nodes := make([]*html.Node, 0)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, c)
	}
	return nodes
}

func compare(t *Tag, n *html.Node) bool {
	if len(t.Attr) == 0 {
		return true
	}
	nTag := node2tag(n)
	for k, v := range t.Attr {
		if _, ok := nTag.Attr[k]; ok {
			for kk, _ := range v {
				if _, ok := nTag.Attr[k][kk]; !ok {
					return false
				}
			}
		} else {
			return false
		}
	}
	return true
}

func node2tag(n *html.Node) *Tag {
	tagname := n.Data
	attrmap := make(map[string]map[string]struct{}, 1)
	for _, at := range n.Attr {
		if _, ok := attrmap[at.Key]; !ok {
			attrmap[at.Key] = make(map[string]struct{}, 1)
		}
		for _, val := range strings.Split(at.Val, " ") {
			attrmap[at.Key][val] = struct{}{}
		}
	}
	return &Tag{tagname, attrmap}
}
