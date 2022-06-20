package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Println(err)
		return
	}
	nums := []int64{1, 3, 5, 68, 9}
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(msg)

	for i := range nums {
		err = binary.Write(conn, binary.BigEndian, nums[i])
		if err != nil {
			log.Println(err)
			return
		}
	}
}
