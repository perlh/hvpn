package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tunvitor/socks"
)

/**
 * @Author: hsm
 * @Email: hsmcool@163.com
 * @Date: 2023/5/16 15:07
 * @Desc:
 */

var config struct {
	Verbose    bool
	UDPTimeout time.Duration
	TCPCork    bool
}

var TunConfig struct {
	ServerIP   string
	ServerIPv6 string
}

type AppConfig struct {
	socks        string
	ltun         string
	stun         string
	server       string
	remoteServer string
	subnet       string
}

var clientMode bool
var serverMode bool

func main() {
	appConfig := AppConfig{}
	flag.BoolVar(&clientMode, "c", false, "client mode")
	flag.BoolVar(&serverMode, "s", false, "server mode")
	flag.StringVar(&appConfig.socks, "socks", "0.0.0.0:1080", "socks -> client ip")
	flag.StringVar(&appConfig.remoteServer, "remote_server", "222.204.52.222:9000", "listen ss server")
	flag.StringVar(&appConfig.server, "server", "0.0.0.0:9000", "client -> server ip")
	flag.StringVar(&appConfig.ltun, "ltun", "0.0.0.0:30000", "local tun ip")
	flag.StringVar(&appConfig.stun, "stun", "222.204.52.222:30000", "remote/server tun ip")
	flag.StringVar(&appConfig.subnet, "subnet", "192.168.123.1/24", "vpn subnet")
	flag.Parse()
	if clientMode == true && serverMode == false {
		socks.UDPEnabled = true
		log.Println("run mode: client")
		log.Println(appConfig.stun)
		go socksLocal(appConfig.socks, appConfig.remoteServer)
		go udpSocksLocal(appConfig.socks, appConfig.remoteServer)
		go StartUDPClient(appConfig.subnet, appConfig.stun, appConfig.ltun)
	} else if serverMode == true && clientMode == false {
		log.Println("run mode: server")
		go tcpRemote(appConfig.server)
		go udpRemote(appConfig.server)
		go StartUDPServer(appConfig.subnet, appConfig.ltun)

	} else {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if clientMode == true {

	} else {

	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

}
