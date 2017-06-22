package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"github.com/garyburd/redigo/redis"
)
const (
	MAX_CONN_NUM = 100
)
func EchoFunc(conn net.Conn,c redis.Conn) {
	defer conn.Close()
	buf := make([]byte, 10)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			println("Error reading:", conn.RemoteAddr(), err.Error())
			return
		}
		if err == nil {
			fmt.Println(string(buf))
			r,err := c.Do("SADD","vehicle1:330101:201706211943","fdakljfdakl")
			if err != nil {
				fmt.Println(err)
			}
			if err == nil {
				fmt.Println(r)
			}
		}
		//fmt.Printf("server %s\n", string(buf))
		//send reply
		_, err = conn.Write(buf)
		if err != nil {
			println("Error send reply:", conn.RemoteAddr(), err.Error())
			return
		}
	}
}
func main() {
	c, err := redis.Dial("tcp", "10.0.30.120:6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	//defer c.Close()

	listener, err := net.Listen("tcp", "10.70.7.181:8090")
	if err != nil {
		fmt.Println("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf("running ...\n")

	var connNum int = 0
	connChan := make(chan net.Conn)
	connChangeChan := make(chan int)

	go func() {
		for connChange := range connChangeChan {
			connNum += connChange
		}
	}()

	go func() {
		for _ = range time.Tick(5e9) {
			fmt.Printf("connNum: %d\n", connNum)
		}
	}()

	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range connChan {
				connChangeChan <- 1
				EchoFunc(conn,c)
				connChangeChan <- -1
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			println("Error accept:", err.Error())
			return
		}
		connChan <- conn
	}
}