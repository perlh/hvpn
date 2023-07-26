package main

/**
 * @Author: hsm
 * @Email: hsmcool@163.com
 * @Date: 2023/7/3 09:42
 * @Desc:
 */

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func redir() {
	// 定义本地监听地址
	localAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:9001")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}

	// 定义目标服务器地址
	remoteAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:1081")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
	//fmt.Println("connect to 127.0.0.1:1081")
	// 建立本地监听
	listener, err := net.ListenTCP("tcp", localAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		os.Exit(1)
	}
	log.Println("建立本地监听:", localAddr)

	// 循环接收连接并处理
	for {
		conn, err := listener.AcceptTCP()
		log.Println(conn.RemoteAddr())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: %s", err.Error())
			continue
		}

		// 建立与目标服务器的连接
		remoteConn, err := net.DialTCP("tcp", nil, remoteAddr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: %s", err.Error())
			conn.Close()
			continue
		}

		// 启动协程处理数据流转发
		go func() {
			defer remoteConn.Close()
			defer conn.Close()
			_, err := io.Copy(conn, remoteConn)
			if err != nil {
				log.Println("111")
				fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			}
		}()

		// 启动协程处理数据流转发
		go func() {
			defer remoteConn.Close()
			defer conn.Close()

			_, err := io.Copy(remoteConn, conn)
			if err != nil {
				log.Println("222")
				fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
			}
		}()
	}
}
