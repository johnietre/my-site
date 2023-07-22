package main

import (
	"net"
)

var (
	conns map[string]net.Conn
)

func listen() {
	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			println(err.Error())
			continue
		}
		go handleServer(conn)
	}
}

func handleServer(conn net.Conn) {
	//
}

func requestPath(request []byte) (path string) {
	for _, b := range request {
		//
	}
}
