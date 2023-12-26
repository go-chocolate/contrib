package jsonutil

import "testing"

func TestAny(t *testing.T) {
	var data = `{"name":"John","data":[{"id":1,"value":12.3,"enable":true},{"id":2,"enable":false}],"content":{"name":""}}`
	var a = new(Any)

	if err := a.UnmarshalJSON([]byte(data)); err != nil {
		t.Error(err)
		return
	}
	t.Log(a)
}
