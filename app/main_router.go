package main

import (
	"github.com/coreos/pkg/capnslog"
	"staticzeng.com/config"
	"staticzeng.com/udp"
)

var (
	log = capnslog.NewPackageLogger("sctele.com/apps/telecom/server", "server.main")
)

func main() {
	cfg := config.LoadConfig()
	udpserver := udp.NewudpServer(cfg.Server.Host, cfg.Server.Port)
	if err := udpserver.Start(); err != nil {
		log.Error("server start error")
	} else {
		defer udpserver.Stop()
		ch := make(chan bool)
		<-ch
	}
}
