package nodefinder

import "golang.org/x/net/html"
import "io"
import "strings"

//Convert path string to []*Elem for further use.
func NewPath(s string) []*Elem {
	return Parse(Lex(s))
}

//Given a path (converted to []*Elem), find all matched nodes and return their pointers in a slice)
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

	var sp int //startpoint for searching
	if Compare(elems[0], root) {
		roots = append(roots, root)
	} else {
		if elems[0].Root {
			return roots
		} else {
			for i, e := range elems {
				if Empty(e) || e.Tag == ".." {
					continue
				}
				sp = i
				find1(elems[i], root, &roots)
				break
			}
		}
	}
	if len(elems) == sp+1 {
		if elems[sp].Nchild != 0 {
			if elems[sp].Nchild <= len(roots) {
				return roots[elems[sp].Nchild-1 : elems[sp].Nchild]
			} else {
				return []*html.Node{}
			}
		} else {
			return roots
		}
	} else {
		if elems[sp].Nchild != 0 {
			if elems[sp].Nchild <= len(roots) {
				find2(elems, sp+1, roots[elems[sp].Nchild-1], &result)
			} else {
				return []*html.Node{}
			}
		} else {
			for _, f := range roots {
				find2(elems, sp+1, f, &result)
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

	if elems[idx].Tag == ".." && p.Parent.Parent != nil {
		if idx == len(elems)-1 {
			*result = append(*result, p.Parent)
			return
		} else {
			find2(elems, idx+1, p.Parent, result)
		}
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

// Check if e *Elem if empty (contains no tag and attributes).
func Empty(e *Elem) bool {
	if e.Tag == "" && len(e.Attr) == 0 && e.Nchild == 0 {
		return true
	} else {
		return false
	}
}

//Compare e *Elem and n *html.Node, if they have the same tag, and the attributes of e (if has any) is a subset of ones in n, then return true, otherwise return false.
func Compare(e *Elem, n *html.Node) bool {
	if n.Type != html.ElementNode {
		return false
	}
	if Empty(e) {
		return true
	}

	if e.Tag == ".." {
		if n.Parent != nil && n.Parent.Type == html.ElementNode {
			return true
		} else {
			return false
		}
	}

	found := true
	if e.Tag == n.Data || e.Tag == "" {
	MATCH1:
		for key, val := range e.Attr {
			for _, attr := range n.Attr {
				if key == attr.Key {
					tmp1 := strings.Split(val, " ")
					tmp2 := strings.Split(attr.Val, " ")
				MATCH2:
					for _, t1 := range tmp1 {
						for _, t2 := range tmp2 {
							if t1 == t2 {
								continue MATCH2
							}
						}
						// if found match in inner loop, here will never reach
						found = false
						break MATCH1
					}
					continue MATCH1
				}
			}
			// if found match in inner loop, here will never reach
			found = false
			break
		}
	} else {
		return false
	}
	return found
}
