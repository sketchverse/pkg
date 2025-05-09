package snowflake

import "testing"

func TestGenerate(t *testing.T) {
	node, err := NewNode(0, true)
	if err != nil {
		t.Error(err)
	}
	id := node.Generate()
	t.Log(id)
}
