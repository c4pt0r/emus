package emus

import (
	"bytes"
	"fmt"
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
	m      map[string]interface{}
}

func newContext(parent *context) *context {
	return &context{
		parent: parent,
		m:      make(map[string]interface{}),
	}
}

func (c *context) put(key string, val interface{}) {
	c.m[key] = val
}

func (c *context) get(key string, recursive bool) (interface{}, bool) {
	cur := c
	if recursive {
		for cur != nil {
			if val, ok := cur.m[key]; ok {
				return val, ok
			}
			cur = c.parent
		}
		return nil, false
	} else {
		val, ok := cur.m[key]
		return val, ok
	}
}

type token struct {
	key      string // {{ name }} <--- name
	typ      tokenType
	body     segment
	escape   bool
	children []*token
}

func newToken(
	typ tokenType,
	key string,
	body segment,
	escape bool,
	children []*token) *token {
	return &token{
		key:      key,
		typ:      typ,
		body:     body,
		escape:   escape,
		children: children,
	}
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
	if val, ok := ctx.get(t.key, true); ok {
		return val.([]byte)
	}
	return nil
}

func (t *token) renderSection(ctx *context) []byte {
	var ret [][]byte

	for _, child := range t.children {
		ret = append(ret, child.render(ctx))
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
