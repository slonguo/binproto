# binproto

`binproto` 为一轻量的的二进制编解码协议（数据交换格式），特点是跨平台跨语言（没有采用Golang自带的gob来编解码）, 简化了 TLV(Type-Length-Value) 规则（这里 type 只为 string 类型）。

`binproto` 主要是用于展示二进制的编解码，不适合性能等要求高的环境。


## 格式


### 公共：KV map的形式

```go
type Params struct {
    fields map[string]string
}
```

经二进制编码后的格式：

```text
总长度|{key1val1总长|key1长|key1|val1长|val1|}|{key2val2总长|key2长|key2|val2长|val2|}| ...
```

### 二进制请求

```go
type ReqMsg struct {
    method string
    params *Params
}
```

经二进制编码后的请求：

```text
总长度|method长|method名|params的二进制编码
```

### 二进制返回

```go
type RespMsg struct {
    code   int32
    msg    string
    data *Params
}
```

经二进制编码后的返回:

```text
总长度|code值|msg长度|msg内容|params的二进制编码
```


### data段为JSON格式的二进制返回

```go
type RespJsonMsg struct {
    code int32
    msg  string
    data string
}
```

其中data为返回的json正文内容

经二进制编码后的返回：

```text
总长度|code值|msg长度|msg内容|data长度|data内容
```
