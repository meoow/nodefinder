package nodefinder

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
)

const (
	SKIP         = "NONE"
	SEP          = "SEP"
	SEPP         = "SEPP"
	TAG          = "TAG"
	CLASS        = "CLASS"
	ID           = "ID"
	NCHILD       = "NCHILD"
	ATTR_START   = "ATTRS_START"
	ATTR_END     = "ATTRS_END"
	ATTR_KEY_VAL = "ATTR_KEY_VAL"
	ATTR_SGL     = "ATTR_SGL"
)

var _rules = [][2]string{
	{`^/{2,}`, SEPP},
	{`^/`, SEP},
	{`^[[:alpha:]][[:alnum:]]*`, TAG},
	{`^\.[[:alpha:]][[:alnum:]_]*`, CLASS},
	{`^#[[:alpha:]][[:alnum:]_]*`, ID},
	{`^:[1-9][0-9]*`, NCHILD},
	{`^\[`, ATTR_START},
	{`^\s+`, SKIP},
}

var _rules_attr = [][2]string{
	{`^([[:alpha:]][[:alnum:]_]*)\s*=\s*([[:alnum:]_]+)`, ATTR_KEY_VAL},
	{`^([[:alpha:]][[:alnum:]_]*)\s*=\s*("[^"]*")`, ATTR_KEY_VAL},
	{`^([[:alpha:]][[:alnum:]_]*)\s*=\s*('[^']*')`, ATTR_KEY_VAL},
	{`^([[:alpha:]][[:alnum:]_]*)`, ATTR_SGL},
	{`^("[^"]+")`, ATTR_SGL},
	{`^('[^']+')`, ATTR_SGL},
	{`^\]`, ATTR_END},
	{`^\s+`, SKIP},
}

type Token struct {
	Text string
	Tag  string
}

type Rule struct {
	Pattern *regexp.Regexp
	Tag     string
}

var Rules []*Rule
var RulesAttr []*Rule

func init() {
	Rules = make([]*Rule, 0, len(_rules))
	RulesAttr = make([]*Rule, 0, len(_rules_attr))
	for _, r := range _rules {
		Rules = append(Rules, &Rule{regexp.MustCompile(r[0]), r[1]})
	}

	for _, r := range _rules_attr {
		RulesAttr = append(RulesAttr, &Rule{regexp.MustCompile(r[0]), r[1]})
	}
}

func Lex(str string) []*Token {
	tokens := make([]*Token, 0, 1)
	pos := 0
	inAttrsLex := false
	strbytes := []byte(str)

	for len(strbytes) > 0 {
		found := false
		if !inAttrsLex {
			for _, rule := range Rules {
				index := rule.Pattern.FindIndex(strbytes)
				if index != nil {
					found = true
					switch rule.Tag {
					case SKIP:
					case ATTR_START:
						inAttrsLex = true
						fallthrough
					case TAG:
						tokens = append(tokens, &Token{string(strbytes[index[0]:index[1]]), rule.Tag})
					case CLASS, ID:
						tokens = append(tokens, &Token{string(strbytes[index[0]+1 : index[1]]), rule.Tag})
					case SEP, SEPP:
						tokens = append(tokens, &Token{"/", rule.Tag})
					case NCHILD:
						tokens = append(tokens, &Token{string(strbytes[index[0]+1 : index[1]]), rule.Tag})

					}
					strbytes = strbytes[index[1]:]
					pos += index[1]
					break
				}
			}
		} else {
			for _, rule := range RulesAttr {
				index := rule.Pattern.FindSubmatchIndex(strbytes)
				if index != nil {
					found = true
					switch rule.Tag {
					case SKIP:
					case ATTR_KEY_VAL:
						key := strbytes[index[2]:index[3]]
						val := strbytes[index[4]:index[5]]
						if bytes.HasPrefix(val, []byte{'\''}) {
							val = bytes.Trim(val, "'")
						} else if bytes.HasPrefix(val, []byte{'"'}) {
							val = bytes.Trim(val, "\"")
						}
						tokens = append(tokens, &Token{string(bytes.Join([][]byte{key, val}, []byte{' '})), rule.Tag})
					case ATTR_SGL:
						key := strbytes[index[0]:index[1]]
						if bytes.HasPrefix(key, []byte{'\''}) {
							key = bytes.Trim(key, "'")
						} else if bytes.HasPrefix(key, []byte{'"'}) {
							key = bytes.Trim(key, "\"")
						}
						tokens = append(tokens, &Token{string(key), rule.Tag})
					case ATTR_END:
						inAttrsLex = false
						key := strbytes[index[0]:index[1]]
						tokens = append(tokens, &Token{string(key), rule.Tag})
					}
					strbytes = strbytes[index[1]:]
					pos += index[1]
					break
				}
			}
		}

		if !found {
			fmt.Println(string(strbytes))
			log.Fatalln("Syntax error at position: ", pos)
		}
	}
	if len(tokens) > 0 {
		switch tokens[len(tokens)-1].Tag {
		case SEP, SEPP:
			tokens = tokens[:len(tokens)-1]
		}
	}
	return tokens
}
