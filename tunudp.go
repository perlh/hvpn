package main

import (
	"github.com/songgao/water"
	"github.com/songgao/water/waterutil"
	"log"
	"net"
	"strings"
	"sync"
	"time"
	"tunvitor/netutil"
	"tunvitor/tun"
)

type tunmap struct {
	sync.RWMutex
	m       map[string]net.UDPAddr
	timeout time.Duration
}

func newTUNmap(timeout time.Duration) *tunmap {
	m := &tunmap{}
	m.m = make(map[string]net.UDPAddr)
	m.timeout = timeout
	return m
}

func (m *tunmap) Get(key string) net.UDPAddr {
	m.RLock()
	defer m.RUnlock()
	return m.m[key]
}

func (m *tunmap) Set(key string, pc *net.UDPAddr) {
	m.Lock()
	defer m.Unlock()
	m.m[key] = *pc
}

func (m *tunmap) Del(key string) net.UDPAddr {
	m.Lock()
	defer m.Unlock()

	pc, ok := m.m[key]
	if ok {
		delete(m.m, key)
		return pc
	}
	return net.UDPAddr{}
}

var nm2 = newTUNmap(config.UDPTimeout)

/**
 * @Author: hsm
 * @Email: hsmcool@163.com
 * @Date: 2023/7/25 16:06
 * @Desc:
 */

func StartUDPClient(CIDR, ServerAddr, LocalAddr string) {
	iface := tun.CreateTun(CIDR)
	serverAddr, err := net.ResolveUDPAddr("udp", ServerAddr)
	if err != nil {
		log.Fatalln("failed to resolve server addr:", err)
	}

	localAddr, err := net.ResolveUDPAddr("udp", LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}

	defer conn.Close()
	parts := strings.Split(CIDR, "/")
	log.Printf("listen tun over udp %v <-> %v <----> %v <-> %v", parts[0], LocalAddr, serverAddr.String(), CIDR)

	go func() {
		buf := make([]byte, 1500)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil || n == 0 {
				continue
			}
			b := buf[:n]
			if !waterutil.IsIPv4(b) {
				continue
			}
			iface.Write(b)
		}
	}()

	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		if !waterutil.IsIPv4(packet) {
			continue
		}
		b := packet[:n]
		conn.WriteToUDP(b, serverAddr)
	}
}

type Forwarder struct {
	localConn *net.UDPConn
	tunmap    tunmap
}

func (f *Forwarder) forward(iface *water.Interface, conn *net.UDPConn) {
	packet := make([]byte, 1500)
	for {
		n, err := iface.Read(packet)
		if err != nil || n == 0 {
			continue
		}
		b := packet[:n]
		if !waterutil.IsIPv4(b) {
			continue
		}

		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		v := nm2.Get(dstAddr)

		//log.Println("找到客户端ip:", v, dstAddr)
		f.localConn.WriteToUDP(b, &v)
	}
}

func StartUDPServer(CIDR, LocalAddr string) {
	// 创建一个tun虚拟网卡
	iface := tun.CreateTun(CIDR)
	//
	localAddr, err := net.ResolveUDPAddr("udp", LocalAddr)
	if err != nil {
		log.Fatalln("failed to get UDP socket:", err)
	}
	// 监听udp
	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		log.Fatalln("failed to listen on UDP socket:", err)
	}
	defer conn.Close()

	log.Printf("udp server started on %v,CIDR is %v", LocalAddr, CIDR)

	forwarder := &Forwarder{localConn: conn}
	go forwarder.forward(iface, conn)

	buf := make([]byte, 1500)

	for {
		// 从udp中读取消息
		n, cliAddr, err := conn.ReadFromUDP(buf)
		if err != nil || n == 0 {
			continue
		}
		b := buf[:n]
		// 判断不是ipv4流量
		if !waterutil.IsIPv4(b) {
			continue
		}
		iface.Write(b)
		srcAddr, dstAddr := netutil.GetAddr(b)
		if srcAddr == "" || dstAddr == "" {
			continue
		}
		v := nm2.Get(srcAddr)
		if v.String() == cliAddr.String() {
			continue
		}
		nm2.Set(srcAddr, cliAddr)
		//log.Println("accept ip connect:", cliAddr.String())
	}
}
