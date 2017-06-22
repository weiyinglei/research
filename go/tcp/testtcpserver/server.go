package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

const (
	MAX_CONN_NUM = 100000
)

func EchoFunc(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 10)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			//println("Error reading:", err.Error())
			return
		}
		//fmt.Printf("server %s\n", string(buf))
		//send reply
		_, err = conn.Write(buf)
		if err != nil {
			//println("Error send reply:", err.Error())
			return
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "10.70.7.181:8088")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("running ...\n")

	var cur_conn_num int = 0
	conn_chan := make(chan net.Conn)
	ch_conn_change := make(chan int)

	go func() {
		for conn_change := range ch_conn_change {
			cur_conn_num += conn_change
		}
	}()

	go func() {
		for _ = range time.Tick(5e9) {
			fmt.Printf("cur conn num: %d\n", cur_conn_num)
		}
	}()

	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range conn_chan {
				ch_conn_change <- 1
				EchoFunc(conn)
				ch_conn_change <- -1
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		conn_chan <- conn
	}
}