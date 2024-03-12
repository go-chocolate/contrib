package jsonutil

import "testing"

func TestAny(t *testing.T) {
	var data = `{"name":"John","data":[{"id":1,"value":12.3,"enable":true},{"id":2,"enable":false}],"content":{"name":""},"count":65536}`
	var a = new(Any)

	if err := a.UnmarshalJSON([]byte(data)); err != nil {
		t.Error(err)
		return
	}
	t.Log(a)
	t.Log(a.Object())
	object, _ := a.Object()
	t.Log(object["name"])
	t.Log(object["data"].Array())
	t.Log(object["count"].Float64())
	t.Log(object["count"].Float64())
	t.Log(a.Type)
	t.Log(a.Value())
}
