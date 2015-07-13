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
	PARENT       = "PARENT"
)

type Token struct {
	Text string
	Tag  string
}

type Rule struct {
	Pattern *regexp.Regexp
	Tag     string
}

var _rec = regexp.MustCompile

var Rules = []*Rule{
	&Rule{_rec(`^/{2,}`), SEPP},
	&Rule{_rec(`^/`), SEP},
	&Rule{_rec(`^\.\.`), PARENT},
	&Rule{_rec(`^[[:alpha:]][[:alnum:]]*`), TAG},
	&Rule{_rec(`^\.[[:alpha:]][[:alnum:]_-]*`), CLASS},
	&Rule{_rec(`^#[[:alpha:]][[:alnum:]_-]*`), ID},
	&Rule{_rec(`^:[1-9][0-9]*`), NCHILD},
	&Rule{_rec(`^\[`), ATTR_START},
	&Rule{_rec(`^\s+`), SKIP},
}

var RulesAttr = []*Rule{
	&Rule{_rec(`^([[:alpha:]][[:alnum:]_]*)\s*=\s*([[:alnum:]_]+)`), ATTR_KEY_VAL},
	&Rule{_rec(`^([[:alpha:]][[:alnum:]_]*)\s*=\s*("[^"]*")`), ATTR_KEY_VAL},
	&Rule{_rec(`^([[:alpha:]][[:alnum:]_]*)\s*=\s*('[^']*')`), ATTR_KEY_VAL},
	&Rule{_rec(`^([[:alpha:]][[:alnum:]_]*)`), ATTR_SGL},
	&Rule{_rec(`^("[^"]+")`), ATTR_SGL},
	&Rule{_rec(`^('[^']+')`), ATTR_SGL},
	&Rule{_rec(`^\]`), ATTR_END},
	&Rule{_rec(`^\s+`), SKIP},
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
					case PARENT:
						tokens = append(tokens, &Token{"..", rule.Tag})

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
