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
	var tokenStack [][]*token
	left := 0
	right := len(tmpl)
	for _, idx := range indexes {
		if left < idx[0] {
			// add literal token
			t := newToken(LITERAL, "",
				segment{left, idx[0] - 1, tmpl[left : idx[0]-1]},
				false,
				nil,
			)
			tokens = append(tokens, t)
		}

		var t *token
		results := reTag.FindSubmatch(tmpl[idx[0]:idx[1]])
		_, prefix, key := results[0], results[1], results[2]

		switch string(prefix) {
		case "":
			{
				t = newToken(VARIABLE, string(key),
					segment{idx[0], idx[1], tmpl[idx[0]:idx[1]]},
					false,
					nil,
				)
				tokens = append(tokens, t)
			}
		case "#":
			{
				t = newToken(SECTION, string(key),
					segment{idx[0], idx[1], tmpl[idx[0]:idx[1]]},
					false,
					nil,
				)
				tokens = append(tokens, t)
				// push to stack
				tokenStack = append(tokenStack, tokens)
				tokens = nil
			}
		case "/":
			{
				top := tokenStack[len(tokenStack)-1]
				sectionToken := top[len(top)-1]

				if sectionToken.key == string(key) {
					sectionToken.children = tokens
				} else {
					log.Warn("section mismatch")
				}

				// pop stack
				tokens = top
				tokenStack = tokenStack[0 : len(tokenStack)-1]
			}

		}
		left = idx[1]
	}
	if left < right {
		t := newToken(LITERAL, "",
			segment{left, right, tmpl[left:right]},
			false,
			nil,
		)
		tokens = append(tokens, t)
	}

	return newToken(ROOT, "", segment{}, false, tokens)
}
