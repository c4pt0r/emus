package emus

import (
	"regexp"

	"github.com/ngaut/log"
)

func delimiterToRegexp(left, right string) *regexp.Regexp {
	return regexp.MustCompile(left + `([#^>&{/!=]?)\s*(.*?)\s*([}=]?)` + right)
}

// returns the token tree
func parse(tmpl []byte) *token {
	// TODO: use default tag
	reTag := delimiterToRegexp("{{", "}}")

	// search tags
	indexes := reTag.FindAllIndex(tmpl, -1)
	// create token
	var tokens []*token
	var ts [][]*token
	left := 0
	right := len(tmpl)
	for _, idx := range indexes {
		if left < idx[0] {
			// add literal token
			t := &token{
				typ:  LITERAL,
				body: segment{left, idx[0] - 1, tmpl[left : idx[0]-1]},
			}
			tokens = append(tokens, t)
		}

		tag := tmpl[idx[0]:idx[1]]
		results := reTag.FindSubmatch(tag)
		_, prefix, key, suffix := results[0], results[1], results[2], results[3]
		seg := segment{idx[0], idx[1], tag}

		var t *token
		switch string(prefix) {
		case "":
			{
				t = &token{
					typ:    VARIABLE,
					key:    string(key),
					suffix: string(suffix),
					body:   seg,
				}
				tokens = append(tokens, t)
			}
		case "#":
			{
				t = &token{
					typ:    SECTION,
					key:    string(key),
					suffix: string(suffix),
					body:   seg,
				}
				tokens = append(tokens, t)
				// push to stack
				ts = append(ts, tokens)
				// reset token inside section
				tokens = nil
			}
		case "/":
			{
				top := ts[len(ts)-1]
				sectionToken := top[len(top)-1]

				if sectionToken.key == string(key) {
					sectionToken.children = tokens
				} else {
					log.Warn("section mismatch")
				}

				// pop stack
				tokens, ts = ts[len(ts)-1], ts[:len(ts)-1]
			}
		case ">":
			{
				t = &token{
					typ:    PARTIAL,
					key:    string(key),
					suffix: string(suffix),
					body:   seg,
				}
				tokens = append(tokens, t)
			}
		}
		left = idx[1]
	}
	// dont forget last literal token
	if left < right {
		t := &token{
			typ:  LITERAL,
			body: segment{left, right, tmpl[left:right]},
		}
		tokens = append(tokens, t)
	}

	return &token{
		typ:      ROOT,
		children: tokens,
	}
}
