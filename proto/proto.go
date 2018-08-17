package proto

import (
	"github.com/coreos/pkg/capnslog"
	"reflect"
	"staticzeng.com/packet"
	"sync"
)

var (
	RegistedProtos = sync.Map{}
	log            = capnslog.NewPackageLogger(reflect.TypeOf(struct{}{}).PkgPath(), "proto.proto")
)

type Proto interface {
	Handle(*packet.Packet) error
	Code() uint8
}

func registeProto(p Proto) {
	code := p.Code()
	RegistedProtos.Store(code, p)
}

func Dispatch(pkt *packet.Packet) {
	val, ok := RegistedProtos.Load(pkt.Header.ProtoType)
	if !ok {
		log.Error("ProtoType ", pkt.Header.ProtoType, " not support!")
		return
	}
	if pro, ok := val.(Proto); ok {
		pro.Handle(pkt)
	}
}
