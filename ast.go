package emus

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ngaut/log"
)

type tokenType int

const (
	LITERAL  tokenType = iota // text
	VARIABLE                  // {{name}}
	SECTION                   // {{#name}} ... {{/name}}
	INVERTED                  // {{^name}} ... {{/name}}
	COMMENT                   // {{! name }}
	PARTIAL                   // {{> name}}
	ROOT
)

type segment struct {
	start, end int
	v          []byte
}

type context struct {
	parent *context
	v      interface{}
}

func newContext(parent *context, v interface{}) *context {
	return &context{
		parent: parent,
		v:      v,
	}
}

func (c *context) get(path string) (interface{}, bool) {
	// TODO: what about array type?
	cur, ok := c.v.(map[string]interface{})
	if !ok {
		return nil, false
	}
	keys := strings.Split(path, ".")
	for _, key := range keys[:len(keys)-1] {
		v, ok := cur[key]
		if ok {
			cur, ok = v.(map[string]interface{})
			if !ok {
				return nil, false
			}
		} else {
			return nil, false
		}
	}
	v, ok := cur[keys[len(keys)-1]]
	return v, ok
}

type token struct {
	key      string // {{ name }} <--- name
	suffix   string
	typ      tokenType
	body     segment
	children []*token
}

func (t *token) String() string {
	return fmt.Sprintf("[%d,%d]\t%d\t%s", t.body.start, t.body.end, t.typ, string(t.body.v))
}

func (t *token) render(ctx *context) []byte {
	switch t.typ {
	case LITERAL:
		return t.renderLiteral(ctx)
	case VARIABLE:
		return t.renderVarible(ctx)
	case SECTION:
		return t.renderSection(ctx)
	case INVERTED:
		return t.renderInverted(ctx)
	case COMMENT:
		return t.renderComment(ctx)
	case PARTIAL:
		return t.renderPartial(ctx)
	case ROOT:
		return t.renderChildren(ctx)
	}
	return nil
}

func (t *token) renderLiteral(ctx *context) []byte {
	return t.body.v
}

func (t *token) renderVarible(ctx *context) []byte {
	if val, ok := ctx.get(t.key); ok {
		switch vv := val.(type) {
		case string:
			return []byte(vv)
		case int64:
			return []byte(fmt.Sprintf("%lld", vv))
		case float64:
			return []byte(fmt.Sprintf("%llf", vv))
		default:
			log.Warn("unsupported type")
		}
	}
	return nil
}

func (t *token) renderSection(ctx *context) []byte {
	val, ok := ctx.get(t.key)
	if !ok {
		log.Warn("no such key")
		return nil
	}
	var ret [][]byte

	lst, ok := val.([]interface{})
	if !ok {
		log.Warn("not array type")
		return nil
	}

	for _, item := range lst {
		// assert val is a list
		c := newContext(ctx, item)
		for _, child := range t.children {
			ret = append(ret, child.render(c))
		}
	}
	return bytes.Join(ret, []byte(""))
}

func (t *token) renderInverted(ctx *context) []byte {
	return nil
}

func (t *token) renderComment(ctx *context) []byte {
	return nil
}

func (t *token) renderPartial(ctx *context) []byte {
	return nil
}

func (t *token) renderChildren(ctx *context) []byte {
	var ret [][]byte

	for _, child := range t.children {
		ret = append(ret, child.render(ctx))
	}
	return bytes.Join(ret, []byte(""))
}
