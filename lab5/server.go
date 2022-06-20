package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/netip"
	"strings"
)

type Message string

var BlackList map[netip.Addr]Message // in-memory storage for blocked ips and appropriate messages

func main() {
	BlackList = make(map[netip.Addr]Message)      // initializing storage
	addr, err := netip.ParseAddr("192.168.0.173") // adding IP-address to blacklist

	if err != nil {
		log.Println(err)
		return
	}
	BlackList[addr] = "YOU ARE BANNED"
	lst, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Println(err)
		return
	}
	defer func(lst net.Listener) {
		err := lst.Close()
		if err != nil {
			log.Println(err)
		}
	}(lst)

	for {
		conn, err := lst.Accept()
		ipPortStr := conn.RemoteAddr().String() //  retrieve ip address string
		parts := strings.Split(ipPortStr, ":")

		if ip, err := netip.ParseAddr(parts[0]); err == nil {
			if message, ok := BlackList[ip]; ok {
				log.Println(ip, "is in BLACKLIST")
				conn.Write(append([]byte(message), '\n')) // sending message of REJECTION with new line to parse
				conn.Close()                              // closing the connection
				return
			}
		}
		if err != nil {
			log.Println(err)
			return
		}
		go HandleConnection(conn) // new thread for handling the request
	}
}

func HandleConnection(conn net.Conn) {
	log.Println("Accepted new request.Client's IP: ", conn.RemoteAddr().String())
	_, err := conn.Write([]byte("OK\n"))
	if err != nil {
		log.Println(err)
		return
	}

	for {
		var val int64
		err := binary.Read(conn, binary.BigEndian, &val)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Println(val)
	}
	log.Println("request handled.")
	conn.Close()
	return
}
