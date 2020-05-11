package proto

import (
	"net"
	"io"
	"fmt"
	"encoding/binary"
	"bytes"
)

func ParseRespMsg(conn net.Conn, maxBody int32) (*RespMsg, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	var msgLen int32 = 0
	err = binary.Read(bytes.NewBuffer(header), binary.BigEndian, &msgLen)
	if err != nil {
		return nil, err
	}

	if msgLen >= maxBody {
		return nil, fmt.Errorf("request body too long! [error:%d >= %d(max)]", msgLen, maxBody)
	}

	bodyLen := msgLen - 4
	msg := make([]byte, bodyLen)
	_, err = io.ReadFull(conn, msg)
	if err != nil {
		return nil, err
	}

	message := bytes.NewBuffer(nil)
	message.Write(header[0:])
	message.Write(msg[0:])

	respMsg := NewRespMsg()
	err = respMsg.Decode(bytes.NewBuffer(message.Bytes()))
	if err != nil {
		return nil, err
	}

	return respMsg, nil
}

func ParseReqMsg(conn net.Conn, maxBody int32) (*ReqMsg, error) {
	header := make([]byte, 4)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return nil, err
	}

	var msgLen int32 = 0
	err = binary.Read(bytes.NewBuffer(header), binary.BigEndian, &msgLen)
	if err != nil {
		return nil, err
	}

	if msgLen >= maxBody {
		return nil, fmt.Errorf("request body too long! [error:%d >= %d(max)]", msgLen, maxBody)
	}

	bodyLen := msgLen - 4
	msg := make([]byte, bodyLen)
	_, err = io.ReadFull(conn, msg)
	if err != nil {
		return nil, err
	}

	message := bytes.NewBuffer(nil)
	message.Write(header[0:])
	message.Write(msg[0:])

	reqMsg := NewReqMsg()
	err = reqMsg.Decode(bytes.NewBuffer(message.Bytes()))
	if err != nil {
		return nil, err
	}

	return reqMsg, nil
}

func SendReqMsg(conn net.Conn, req *ReqMsg) (int, error) {
	buf := new(bytes.Buffer)
	if buf == nil {
		return 0, fmt.Errorf("new buffer failed! nil buf")
	}
	enErr := req.Encode(buf)
	if enErr != nil {
		return 0, fmt.Errorf("encode reqMsg:%s error:%s", req.Dump(), enErr)
	}
	sendByte, err := conn.Write(buf.Bytes())

	return sendByte, err
}

func SendRespMsg(conn net.Conn, resp *RespMsg) (int, error) {
	buf := new(bytes.Buffer)
	if buf == nil {
		return 0, fmt.Errorf("new buffer failed! nil buf")
	}
	enErr := resp.Encode(buf)
	if enErr != nil {
		return 0, fmt.Errorf("encode respMsg:%s error:%s", resp.Dump(), enErr)
	}
	sendByte, err := conn.Write(buf.Bytes())

	return sendByte, err
}