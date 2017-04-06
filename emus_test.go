package emus

import (
	"testing"

	"github.com/ngaut/log"
	"github.com/stretchr/testify/assert"
)

func TestCompileTmpl(t *testing.T) {
	assert := assert.New(t)

	r := parse([]byte(`hello {{#test}} ddd {{ aa }} {{#test2}} fff {{ aaa }} {{/test2}} {{/test}} world`))

	assert.Equal(len(r.children), 3)
}

func TestContextLookup(t *testing.T) {
	assert := assert.New(t)

	c := newContext(nil, map[string]interface{}{
		"1": map[string]interface{}{
			"xx": "yy",
			"2": map[string]interface{}{
				"3": "fuck",
				"4": "ddd",
			},
		},
	})

	v, _ := c.get("1.2.3")
	assert.Equal(v, "fuck")
	v, _ = c.get("1.xx")
	assert.Equal(v, "yy")
	v, _ = c.get("1.x")
	assert.Equal(v, nil)
}

func TestRender(t *testing.T) {
	tmpl := parse([]byte(`
	hello {{name}},
	{{#items}}
	item 
	{{/items}}
	`))
	out := tmpl.render(newContext(nil, map[string]interface{}{
		"name": "dongxu",
		"items": []interface{}{
			1, 2, 3,
		},
	}))
	log.Info(string(out))
}
