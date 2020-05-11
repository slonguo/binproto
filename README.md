# `binproto`

`binproto` is a simple protocol implementation of binary encoding and decoding(data interchange format). It is cross-platform, cross-language(not using Golang gob encoding/decoding) and simplifies the Type-Length-Value(TLV) by using only string type.

`binproto` is only a demo for illustrating binary encoding/decoding. Use sophisticated interchange formats in production whenever possible.

## Format

Params(KV):

```go
type Params struct {
    fields map[string]string
}
```

Params Encoding:

```text
total length|{k1v1 len|k1 len|k1|v1 len|v1|}|{k2v2 len|k2 len|k2|v2 len|v2|}| ...
```

Request:

```go
type ReqMsg struct {
    method string
    params *Params
}
```

Request Encoding:

```text
total length|method len|method|params encoding
```

Response:

```go
type RespMsg struct {
    code   int32
    msg    string
    data *Params
}
```

Response Encoding:

```text
total length|code|msg len|msg|params encoding
```

JSON Response(data is of JSON type):

```go
type RespJsonMsg struct {
    code int32
    msg  string
    data string
}
```

JSON Response Encoding:

```text
total length|code|msg len|msg|data len|data
```