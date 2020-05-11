package proto

import (
	"io"
	"encoding/binary"
	"fmt"
)

func GetLength(buf io.Reader) (error, int32) {
	var s int32 = 0
	err := binary.Read(buf, binary.BigEndian, &s)
	if err != nil {
		return err, 0
	}
	return nil, s
}

type Params struct {
	fields map[string]string
}

func (v *Params) Init() *Params {
	v.fields = make(map[string]string)
	return v
}

func (v *Params) Set(key string, value string) {
	v.fields[key] = value
}

func (v *Params) Del(key string) {
	delete(v.fields, key)
}

func (v *Params) Get(key string) string {
	return v.fields[key]
}

func (v *Params) GetOK(key string) (string, bool) {
	res, ok := v.fields[key]
	return res, ok
}

func (v *Params) GetContent() map[string]string {
	return v.fields
}

func (v *Params) Size() int32 {
	var s int32 = 0
	for k, value := range v.fields {
		s += 12
		s = s + int32(len(k))
		s = s + int32(len(value))
	}
	return s
}

func (v *Params) Encode(buf io.Writer) error {
	size := v.Size()
	if err := binary.Write(buf, binary.BigEndian, size); err != nil {
		return err
	}
	for key, value := range v.fields {
		fieldLen := len(key) + len(value) + 8
		if err := binary.Write(buf, binary.BigEndian, int32(fieldLen)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(key))); err != nil {
			return err
		}
		if _, err := buf.Write([]byte(key)); err != nil {
			return err
		}
		if err := binary.Write(buf, binary.BigEndian, int32(len(value))); err != nil {
			return err
		}
		if _, err := buf.Write([]byte(value)); err != nil {
			return err
		}
	}
	return nil
}

func (v *Params) Decode(buf io.Reader) (error, int32) {
	var ParamsLen int32
	var err error
	if err, ParamsLen = GetLength(buf); err != nil {
		return err, ParamsLen
	}
	remain := ParamsLen
	for remain > 0 {
		var fieldLen, keyLen, valueLen int32
		if err, fieldLen = GetLength(buf); err != nil {
			return err, ParamsLen - remain
		}
		if err, keyLen = GetLength(buf); err != nil {
			return err, ParamsLen - remain
		}
		key := make([]byte, keyLen)
		if _, err = io.ReadFull(buf, key); err != nil {
			return err, ParamsLen - remain
		}
		if err, valueLen = GetLength(buf); err != nil {
			return err, ParamsLen - remain
		}

		value := make([]byte, valueLen)
		if _, err := io.ReadFull(buf, value); err != nil {
			return err, ParamsLen - remain
		}
		v.fields[string(key)] = string(value)
		remain -= fieldLen + 4
	}
	return nil, ParamsLen + 4
}

func (v *Params) Dump() string {
	var result string
	for k, value := range v.fields {
		result += fmt.Sprintf("%s:%s;", k, value)
	}
	return result
}

type ReqMsg struct {
	method string
	params *Params
}

func (v *ReqMsg) Init() *ReqMsg {
	v.params = NewParams()
	return v
}

func (v *ReqMsg) SetMethod(method string) *ReqMsg {
	v.method = method
	return v
}

func (v *ReqMsg) SetParams(params *Params) *ReqMsg {
	v.params = params
	return v
}

func (v *ReqMsg) SetParamsMap(kvs *map[string]string) *ReqMsg {
	if kvs == nil {
		return v
	}
	for key, val := range *kvs {
		v.params.Set(key, val)
	}
	return v
}

func (v *ReqMsg) SetMsg(method string, r *Params) *ReqMsg {
	v.method = method
	v.params = r
	return v
}

func (v *ReqMsg) SetMsgWithMap(method string, kvs *map[string]string) *ReqMsg {
	return v.SetMethod(method).SetParamsMap(kvs)
}

func (v *ReqMsg) GetMsg() (string, map[string]string) {
	return v.method, v.params.GetContent()
}

func (v *ReqMsg) GetMethod() string {
	return v.method
}

func (v *ReqMsg) GetParams() *Params {
	return v.params
}

func (v *ReqMsg) GetParamVal(key string) string {
	return v.params.Get(key)
}

func (v *ReqMsg) GetParamsMap() map[string]string {
	return v.params.GetContent()
}

func (v *ReqMsg) Size() int32 {
	return 12 + int32(len(v.method)) + v.params.Size()
}

func (v *ReqMsg) Encode(buf io.Writer) error {
	msgLen := v.Size()
	if err := binary.Write(buf, binary.BigEndian, msgLen); err != nil {
		return err
	}

	methodLen := len(v.method)
	if err := binary.Write(buf, binary.BigEndian, int32(methodLen)); err != nil {
		return err
	}

	if _, err := buf.Write([]byte(v.method)); err != nil {
		return err
	}

	if err := v.params.Encode(buf); err != nil {
		return err
	}

	return nil
}

func (v *ReqMsg) Decode(buf io.Reader) error {
	err, _ := GetLength(buf)
	if err != nil {
		return err
	}

	var methodLen int32
	if err, methodLen = GetLength(buf); err != nil {
		return err
	}

	method := make([]byte, methodLen)
	if _, err = io.ReadFull(buf, method); err != nil {
		return err
	}

	v.method = string(method)

	params := new(Params).Init()
	if err, _ := params.Decode(buf); err != nil {
		return err
	}

	v.params = params
	return nil
}

func (v *ReqMsg) Dump() string {
	return fmt.Sprintf("%s|%s", v.method, v.params.Dump())
}

type RespMsg struct {
	code   int32
	msg    string
	params *Params
}

func (v *RespMsg) Init() *RespMsg {
	v.params = NewParams()
	return v
}

func (v *RespMsg) SetCodeAndMsg(code int32, msg string) *RespMsg {
	v.code = code
	v.msg = msg
	return v
}

func (v *RespMsg) SetParams(params *Params) *RespMsg {
	v.params = params
	return v
}

func (v *RespMsg) SetParamsMap(kvs *map[string]string) *RespMsg {
	if kvs == nil {
		return v
	}
	for key, val := range *kvs {
		v.params.Set(key, val)
	}
	return v
}

func (v *RespMsg) SetBodyWithMap(code int32, msg string, kvs *map[string]string) *RespMsg {
	v.SetCodeAndMsg(code, msg)
	v.SetParamsMap(kvs)
	return v
}

func (v *RespMsg) SetBody(code int32, msg string, r *Params) *RespMsg {
	v.code = code
	v.msg = msg
	v.SetParams(r)
	return v
}

func (v *RespMsg) GetBody() (int32, string, map[string]string) {
	return v.code, v.msg, v.params.GetContent()
}

func (v *RespMsg) GetCodeAndMsg() (int32, string) {
	return v.code, v.msg
}

func (v *RespMsg) GetParams() *Params {
	return v.params
}

func (v *RespMsg) GetParamVal(key string) string {
	return v.params.Get(key)
}

func (v *RespMsg) GetCode() int32 {
	return v.code
}

func (v *RespMsg) GetMsg() string {
	return v.msg
}

func (v *RespMsg) GetParamsMap() map[string]string {
	return v.params.GetContent()
}

func (v *RespMsg) Size() int32 {
	return 16 + int32(len(v.msg)) + v.params.Size()
}

func (v *RespMsg) Encode(buf io.Writer) error {
	totalLen := v.Size()
	if err := binary.Write(buf, binary.BigEndian, totalLen); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, v.code); err != nil {
		return err
	}

	msgLen := len(v.msg)
	if err := binary.Write(buf, binary.BigEndian, int32(msgLen)); err != nil {
		return err
	}

	if _, err := buf.Write([]byte(v.msg)); err != nil {
		return err
	}

	if err := v.params.Encode(buf); err != nil {
		return err
	}

	return nil
}

func (v *RespMsg) Decode(buf io.Reader) error {
	err, _ := GetLength(buf)
	if err != nil {
		return err
	}

	var code int32
	if err, code = GetLength(buf); err != nil {
		return err
	}

	v.code = code

	var msgLen int32
	if err, msgLen = GetLength(buf); err != nil {
		return err
	}
	msg := make([]byte, msgLen)
	if _, err = io.ReadFull(buf, msg); err != nil {
		return err
	}

	v.msg = string(msg)

	params := new(Params).Init()
	if err, _ := params.Decode(buf); err != nil {
		return err
	}

	v.params = params
	return nil
}

func (v *RespMsg) Dump() string {
	return fmt.Sprintf("[%d|%s|%s]", v.code, v.msg, v.params.Dump())
}

type RespJsonMsg struct {
	code int32
	msg  string
	data string
}

func (v *RespJsonMsg) Init() *RespJsonMsg {
	return &RespJsonMsg{}
}

func (v *RespJsonMsg) SetCodeAndMsg(code int32, msg string) *RespJsonMsg {
	v.code = code
	v.msg = msg
	return v
}

func (v *RespJsonMsg) SetData(data string) *RespJsonMsg {
	v.data = data
	return v
}

func (v *RespJsonMsg) SetDataByBytes(data []byte) *RespJsonMsg {
	v.data = string(data)
	return v
}

func (v *RespJsonMsg) Set(code int32, msg string, data string) *RespJsonMsg {
	v.code = code
	v.msg = msg
	v.data = data
	return v
}

func (v *RespJsonMsg) GetBody() (int32, string, string) {
	return v.code, v.msg, v.data
}

func (v *RespJsonMsg) GetCodeAndMsg() (int32, string) {
	return v.code, v.msg
}

func (v *RespJsonMsg) Size() int32 {
	return 16 + int32(len(v.msg)) + int32(len(v.data))
}

func (v *RespJsonMsg) Encode(buf io.Writer) error {
	totalLen := v.Size()
	if err := binary.Write(buf, binary.BigEndian, totalLen); err != nil {
		return err
	}

	if err := binary.Write(buf, binary.BigEndian, v.code); err != nil {
		return err
	}

	msgLen := len(v.msg)
	if err := binary.Write(buf, binary.BigEndian, int32(msgLen)); err != nil {
		return err
	}

	if _, err := buf.Write([]byte(v.msg)); err != nil {
		return err
	}

	dataLen := len(v.data)
	if err := binary.Write(buf, binary.BigEndian, int32(dataLen)); err != nil {
		return err
	}

	if _, err := buf.Write([]byte(v.data)); err != nil {
		return err
	}

	return nil
}

func (v *RespJsonMsg) Decode(buf io.Reader) error {
	err, _ := GetLength(buf)
	if err != nil {
		return err
	}

	var code int32
	if err, code = GetLength(buf); err != nil {
		return err
	}

	v.code = code

	var msgLen int32
	if err, msgLen = GetLength(buf); err != nil {
		return err
	}
	msg := make([]byte, msgLen)
	if _, err = io.ReadFull(buf, msg); err != nil {
		return err
	}

	v.msg = string(msg)

	var dataLen int32
	if err, dataLen = GetLength(buf); err != nil {
		return err
	}
	data := make([]byte, dataLen)
	if _, err = io.ReadFull(buf, data); err != nil {
		return err
	}

	v.data = string(data)

	return nil
}

func (v *RespJsonMsg) Dump() string {
	return fmt.Sprintf("[%d|%s|%s]", v.code, v.msg, v.data)
}

func NewParams() *Params {
	return new(Params).Init()
}

func NewReqMsg() *ReqMsg {
	return new(ReqMsg).Init()
}

func NewRespMsg() *RespMsg {
	return new(RespMsg).Init()
}

func NewRespJsonMsg() *RespJsonMsg {
	return new(RespJsonMsg).Init()
}

func CreateRespMsg(code int32, msg string, kvs *map[string]string) *RespMsg {
	respMsg := NewRespMsg()
	respMsg.SetBodyWithMap(code, msg, kvs)
	return respMsg
}
