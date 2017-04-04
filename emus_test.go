package emus

import (
	"fmt"
	"testing"
)

func TestCompileTmpl(t *testing.T) {
	r := parse([]byte(`hello {{#test}} ddd {{ aa }} {{#test2}} fff {{ aaa }} {{/test2}} {{/test}} world`))
	fmt.Println(r.render(nil))
}
