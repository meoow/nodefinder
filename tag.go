package nodefinder

import "regexp"
import "strings"

type Tag struct {
	Tag  string
	Attr map[string]map[string]struct{}
}

var toks = regexp.MustCompile(`(?:^[^.#\[]+|[.#]([^.#\[]+)|\[\s*([^\]]+?)\s*\])`)
var atk = regexp.MustCompile(`(?:([^,[=]+)=(?:"([^"]*?[^\\])"|'([^']*?)'|([^,\]]+?)))\s*[,\]]`)

func TagParser(path string) []*Tag {
	tokens := Split(path)
	tags := make([]*Tag, 0, len(tokens))
	for _, t := range tokens {
		tt := tagParser(t)
		if tt != nil {
			tags = append(tags, tt)
		}
	}
	return tags
}

func tagParser(element string) *Tag {
	tokens := toks.FindAllStringSubmatch(element, -1)
	if len(tokens) == 0 {
		return nil
	}
	tag := &Tag{"", make(map[string]map[string]struct{}, 1)}
	for _, t := range tokens {
		switch t[0][0] {
		case '.':
			if _, ok := tag.Attr["class"]; !ok {
				tag.Attr["class"] = make(map[string]struct{}, 1)
			}
			tag.Attr["class"][t[1]] = struct{}{}
		case '#':
			if _, ok := tag.Attr["id"]; !ok {
				tag.Attr["id"] = make(map[string]struct{}, 1)
			}
			tag.Attr["id"][t[1]] = struct{}{}
		case '[':
			attrtokens := atk.FindAllStringSubmatch(t[0], -1)
			for _, a := range attrtokens {
				if _, ok := tag.Attr[a[1]]; !ok {
					tag.Attr[a[1]] = make(map[string]struct{}, 1)
				}
				val := strings.Join(a[2:5], "")
				for _, v := range strings.Split(val, " ") {
					tag.Attr[a[1]][v] = struct{}{}
				}
			}
		default:
			if tag.Tag == "" {
				tag.Tag = t[0]
			}
		}
	}
	return tag
}

func Split(p string) []string {
	toks := make([]string, 0, 1)
	seps := make([]byte, 0, 1)
	expect := byte('0')
	sreader := strings.NewReader(p)
	for {
		b, e := sreader.ReadByte()
		if e != nil {
			if len(seps) > 0 {
				toks = append(toks, string(seps))
			}
			break
		}
		switch b {
		case '[':
			seps = append(seps, b)
			expect = ']'
		case ']':
			seps = append(seps, b)
			if expect == ']' {
				expect = '/'
			}
		case '/':
			if expect == '/' {
				toks = append(toks, string(seps))
				seps = make([]byte, 0, 1)
			} else {
				seps = append(seps, b)
			}
		case ' ', '\t':
			if expect != '/' {
				seps = append(seps, b)
			}
		default:
			if expect == '0' {
				expect = '/'
			}
			seps = append(seps, b)
		}
	}
	return toks
}
