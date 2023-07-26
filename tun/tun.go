package tun

/**
 * @Author: hsm
 * @Email: hsmcool@163.com
 * @Date: 2023/7/23 17:58
 * @Desc:
 */

import (
	"github.com/songgao/water"
	"log"
	"net"
	"os"
	"os/exec"
	"runtime"
)

func CreateTun(cidr string) (iface *water.Interface) {
	c := water.Config{DeviceType: water.TUN}
	iface, err := water.New(c)
	if err != nil {
		log.Fatalln("failed to allocate vpn interface:", err)
	}

	log.Println("interface allocated:", iface.Name())

	ConfigTun(cidr, iface)

	return iface
}

func ConfigTun(cidr string, iface *water.Interface) {
	os := runtime.GOOS
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Panicf("error cidr %v", cidr)
	}

	if os == "linux" {
		execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "mtu", "1500")
		execCmd("/sbin/ip", "addr", "add", cidr, "dev", iface.Name())
		execCmd("/sbin/ip", "link", "set", "dev", iface.Name(), "up")
	} else if os == "darwin" {
		minIp := ipNet.IP.To4()
		minIp[3]++
		execCmd("ifconfig", iface.Name(), "inet", ip.String(), minIp.String(), "up")
	} else if os == "windows" {
		log.Printf("please install openvpn client,see this link:%v", "https://github.com/OpenVPN/openvpn")
		log.Printf("open new cmd and enter:netsh interface ip set address name=\"%v\" source=static addr=%v mask=%v gateway=none", iface.Name(), ip.String(), ipNet.Mask.String())
	} else {
		log.Printf("not support os:%v", os)
	}
}

func execCmd(c string, args ...string) {
	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		log.Fatalln("failed to exec /sbin/ip error:", err)
	}
}
