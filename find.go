package nodefinder

import "code.google.com/p/go.net/html"
import "io"
import "strings"
import "sort"

//import "fmt"

type _string []string

func (s _string) Len() int           { return len(s) }
func (s _string) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s _string) Less(i, j int) bool { return s[i] < s[j] }

func NewPath(s string) []*Elem {
	return Parse(Lex(s))
}

func Find(elems []*Elem, hf io.Reader) ([]*html.Node, error) {
	if len(elems) == 0 {
		return []*html.Node{}, nil
	}

	root, err := html.Parse(hf)
	if err != nil {
		return []*html.Node{}, err
	}

	return FindByNode(elems, root), nil
}

func FindByNode(elems []*Elem, root *html.Node) []*html.Node {

	roots := make([]*html.Node, 0, 1)
	result := make([]*html.Node, 0, 1)

	if Compare(elems[0], root) {
		roots = append(roots, root)
	} else {
		if elems[0].Root {
			return roots
		} else {
			find1(elems[0], root, &roots)
		}
	}
	if len(elems) == 1 {
		if elems[0].Nchild != 0 {
			if elems[0].Nchild <= len(roots) {
				return roots[elems[0].Nchild-1 : elems[0].Nchild]
			} else {
				return []*html.Node{}
			}
		} else {
			return roots
		}
	} else {
		if elems[0].Nchild != 0 {
			if elems[0].Nchild <= len(roots) {
				find2(elems, 1, roots[elems[0].Nchild-1], &result)
			} else {
				return []*html.Node{}
			}
		} else {
			for _, f := range roots {
				find2(elems, 1, f, &result)
			}
		}
	}
	return result
}

func find1(elem *Elem, p *html.Node, found *[]*html.Node) {
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if Compare(elem, c) {
			*found = append(*found, c)
		}
		find1(elem, c, found)
	}
}

func find2(elems []*Elem, idx int, p *html.Node, result *[]*html.Node) {
	if idx >= len(elems) {
		return
	}
	match_count := 0
	for c := p.FirstChild; c != nil; c = c.NextSibling {
		if c.Type != html.ElementNode {
			continue
		}
		if Empty(elems[idx]) {
			if Compare(elems[idx+1], c) {
				match_count++
				if elems[idx+1].Nchild != 0 {
					if elems[idx+1].Nchild == match_count {
						if len(elems)-1 == idx+1 {
							*result = append(*result, c)
						} else {
							find2(elems, idx+2, c, result)
						}
						break
					}
				} else {
					if len(elems)-1 == idx+1 {
						*result = append(*result, c)
					} else {
						find2(elems, idx+2, c, result)
					}
				}
			} else {
				find2(elems, idx, c, result)
			}
		} else if Compare(elems[idx], c) {
			match_count++
			if elems[idx].Nchild != 0 {
				if elems[idx].Nchild == match_count {
					if len(elems)-1 == idx {
						*result = append(*result, c)
					} else {
						find2(elems, idx+1, c, result)
					}
					break
				}
			} else {
				if len(elems)-1 == idx {
					*result = append(*result, c)
				} else {
					find2(elems, idx+1, c, result)
				}
			}
		}
	}
}

func Empty(e *Elem) bool {
	if e.Tag == "" && len(e.Attr) == 0 {
		return true
	} else {
		return false
	}
}

func Compare(e *Elem, n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if Empty(e) {
		return true
	}

	matched_key_count := 0
	if e.Tag == n.Data {
		for key, val := range e.Attr {
			for _, attr := range n.Attr {
				if key == attr.Key {
					tmp1 := strings.Split(val, " ")
					tmp2 := strings.Split(attr.Val, " ")
					sort.Sort(_string(tmp1))
					sort.Sort(_string(tmp2))
					tmps1 := "\x1f" + strings.Join(tmp1, "\x1f") + "\x1f"
					tmps2 := "\x1f" + strings.Join(tmp2, "\x1f") + "\x1f"
					if strings.Contains(tmps2, tmps1) {
						matched_key_count++
					}
				}
			}
		}
	} else {
		return false
	}
	if matched_key_count == len(e.Attr) {
		return true
	} else {
		return false
	}
	panic("")
}
