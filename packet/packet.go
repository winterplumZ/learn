package packet

import (
	"bytes"
	"encoding/binary"
	"net"
)

var (
	HEAD_LENGTH = 4
	MAX_LENGTH  = 10240
)

type PacketHeader struct {
	Magic        uint8
	ProtoType    uint8
	PacketLength uint16
}

type Packet struct {
	Header  *PacketHeader
	NetAddr net.Addr
	Body    []byte
}

func DecodeHeader(header []byte) (p *PacketHeader, err error) {
	defer func() {
		e := recover()
		if e == nil {
			return
		}
		if panicErr, ok := e.(error); ok {
			err = panicErr
		} else {
			panic(e)
		}
	}()
	reader := bytes.NewReader(header)
	p = &PacketHeader{}
	binary.Read(reader, binary.BigEndian, &p.Magic)
	binary.Read(reader, binary.BigEndian, &p.ProtoType)
	binary.Read(reader, binary.BigEndian, &p.PacketLength)
	return
}
