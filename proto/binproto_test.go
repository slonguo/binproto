package proto

import (
	"testing"
	"bytes"
)

func TestParams(t *testing.T) {
	p := NewParams()
	p.Set("h", "j")
	p.Set("l", "m")
	buf := new(bytes.Buffer)
	p.Encode(buf)

	t.Log(p.Size())
}

func TestAdd(t *testing.T) {
	v := NewReqMsg()
	buf := new(bytes.Buffer)
	v.Encode(buf)

	if len(buf.Bytes()) == 12 {
		t.Log("OK")
	} else {
		t.Log("test failed, buf bytes:", len(buf.Bytes()))
	}
}

func TestReqMsg(t *testing.T) {
	v := NewReqMsg()
	v.GetParams().Set("hello", "world")
	v.GetParams().Set("hello", "")

	buf := new(bytes.Buffer)
	v.Encode(buf)
	t.Log(v.Size())

	x := NewReqMsg()
	x.Decode(bytes.NewBuffer(buf.Bytes()))

	if len(x.method) == 0 {
		t.Log("cmd empty ok")
	}

	val := x.GetParamVal("hello")
	newVal := x.GetParamVal("world")
	if val == "world" {
		t.Log("key test ok")
	} else {
		t.Log("value for hello is:", val)
	}

	t.Log(x.Dump())
	t.Log(newVal)
	t.Log(x.GetParamsMap())
	t.Log(x.Size())
}

func TestRespMsg(t *testing.T) {
	v := NewRespMsg()
	v.SetCodeAndMsg(200, "")
	v.GetParams().Set("manga", "sutra")
	v.GetParams().Set("json", "xml")

	t.Log(v.Dump())

	buf := new(bytes.Buffer)
	v.Encode(buf)
	t.Log(v.Size())

	x := NewRespMsg()
	x.Decode(bytes.NewBuffer(buf.Bytes()))

	t.Log(x.GetParamVal("manga"))
	t.Log(x.Dump())
	t.Log(x.code + 1)
}

func TestRespJsonMsg(t *testing.T) {
	v := NewRespJsonMsg()
	v.SetCodeAndMsg(200, "12345")
	v.SetData("{\"key\": 1}")

	t.Log(v.Dump())

	buf := new(bytes.Buffer)
	v.Encode(buf)
	t.Log(v.Size())

	x := NewRespJsonMsg()
	x.Decode(bytes.NewBuffer(buf.Bytes()))
	t.Log(x.Dump())
}