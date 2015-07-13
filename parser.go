package nodefinder

import "strings"
import "log"
import "strconv"

type Elem struct {
	Tag    string
	Attr   map[string]string
	Root   bool
	Nchild int
}

func Parse(toks []*Token) []*Elem {
	elems := make([]*Elem, 0, 1)

	if len(toks) == 0 {
		return elems
	}

	var lastElem *Elem
	lastElem = &Elem{Attr: make(map[string]string, 0)}
	elems = append(elems, lastElem)
	for idx, tok := range toks {
		switch tok.Tag {
		case SEP:
			if idx == 0 {
				lastElem.Root = true
			} else {
				lastElem = &Elem{Attr: make(map[string]string, 0)}
				elems = append(elems, lastElem)
			}
		case SEPP:
			if idx == 0 {
				// no nothing
			} else {
				elems = append(elems, &Elem{})
				lastElem = &Elem{Attr: make(map[string]string, 0)}
				elems = append(elems, lastElem)
			}
		case PARENT:
			fallthrough
		case TAG:
			lastElem.Tag = tok.Text
		case CLASS:
			lastElem.Attr["class"] = tok.Text
		case ID:
			lastElem.Attr["id"] = tok.Text
		case ATTR_KEY_VAL:
			kv := strings.SplitN(tok.Text, " ", 2)
			lastElem.Attr[kv[0]] = kv[1]
		case ATTR_SGL:
			lastElem.Attr[tok.Text] = ""
		case NCHILD:
			tmp, err := strconv.ParseInt(tok.Text, 10, 32)
			lastElem.Nchild = int(tmp)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return elems
}
