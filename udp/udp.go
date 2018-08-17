package udp

import (
	"fmt"
	"github.com/coreos/pkg/capnslog"
	"net"
	"reflect"
	"staticzeng.com/packet"
	"staticzeng.com/proto"
)

var (
	UDP_NETWORK = "udp4"
	log         = capnslog.NewPackageLogger(reflect.TypeOf(struct{}{}).PkgPath(), "udp.udpserver")
)

type udpServer struct {
	host    string
	port    string
	brkChan chan bool
	conn    *net.UDPConn
}

func NewudpServer(host string, port string) *udpServer {
	return &udpServer{
		host: host,
		port: port,
	}
}

func (u *udpServer) Stop() error {
	u.brkChan <- true
	log.Debug("udpServer stop success")
	return nil
}

func (u *udpServer) Start() error {
	if addr, err := net.ResolveUDPAddr(UDP_NETWORK, fmt.Sprintf("%s:%s", u.host, u.port)); err != nil {
		log.Error("ResolveUDPAddr error :", err)
		return err
	} else if conn, err := net.ListenUDP(UDP_NETWORK, addr); err != nil {
		log.Error("ListenUDP error : ", err)
		return err
	} else {
		log.Info("Udpserver listen at :", addr)
		u.conn = conn
		// ready to read
		go u.receiveDataPacket()
		log.Info("Udpserver Start Success")
		return nil
	}
}

func (u *udpServer) receiveDataPacket() {
	defer u.conn.Close()
	for {
		select {
		case <-u.brkChan:
			return
		default:
			packetBytes := make([]byte, packet.MAX_LENGTH)
			if rn, remoteAddr, err := u.conn.ReadFromUDP(packetBytes); err != nil {
				//				log.Info("received rn : ", rn, " content : ", packetBytes)
				log.Error("Read from udp error : ", err)
			} else if rn < packet.HEAD_LENGTH {
				//				log.Info("received rn : ", rn, " content : ", packetBytes)
				log.Error("Read Bytes less than head length")
			} else if header, err := packet.DecodeHeader(packetBytes[:packet.HEAD_LENGTH]); err != nil {
				//				log.Info("received rn : ", rn, " content : ", packetBytes)
				log.Error("Decode header error : ", err)
			} else if rn != packet.HEAD_LENGTH+int(header.PacketLength) {
				//				log.Info("received rn : ", rn, " content : ", packetBytes)
				log.Error("rn != HEAD_LENGTH + PacketLength")
			} else {
				log.Info("received ", rn, " bytes with addr ", remoteAddr.String())
				//				log.Info("received content: ", packetBytes)
				go proto.Dispatch(&packet.Packet{
					Header:  header,
					NetAddr: remoteAddr,
					Body:    packetBytes[packet.HEAD_LENGTH:rn],
				})
			}
		}
	}
}
