package main

import (
	_ "bytes"
	"fmt"
	"net"
	"net/url"
	_ "os"
	"strconv"
	"time"
)

var (
	//host = "tcp://127.0.0.1:8088"
	//host = "tcp://10.0.30.120:8090"
	//host = "tcp://10.0.30.124:8090"
	host = "tcp://10.70.7.181:8090"

	DefalutTimeout = 25 * time.Second
	MaxClient      = 100
	clientNum      = 0
	msgNum = 0
)

func sockConn(daemon string, timeout time.Duration) (net.Conn, error) {
	daemonURL, err := url.Parse(daemon)
	//fmt.Printf("daemon url %v %v \n", daemonURL, daemonURL.Scheme)
	if err != nil {
		return nil, fmt.Errorf("could not parse url %q: %v", daemon, err)
	}

	var c net.Conn
	switch daemonURL.Scheme {
	case "unix":
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Path, timeout)
	case "tcp":
		return net.DialTimeout(daemonURL.Scheme, daemonURL.Host, timeout)
	default:
		return c, fmt.Errorf("unknown scheme %v (%s)", daemonURL.Scheme, daemon)
	}
}

func sendData(conn net.Conn, n int, ch chan int) {
	buf := make([]byte, 10)
	for {
		_, err := conn.Write([]byte(strconv.Itoa(n)))
		if err != nil {
			fmt.Printf("Error reading:%s\n", err.Error())
			clientNum--
			return
		}
		//send reply
		_, err = conn.Read(buf)
		fmt.Printf("client %v\n", string(buf))
		if err != nil {
			fmt.Printf("Error send reply:%s\n", err.Error())
			clientNum--
			return
		}
		//time.Sleep(1 * time.Second)
		ch <- 1
		buf = buf[0:]
	}
}

func connectServer(ch chan int) {
	for i := 1; i <= MaxClient; i++ {
		conn, err := sockConn(host, DefalutTimeout)
		if err != nil {
			fmt.Printf("connect error:%s\n", err)
		} else {
			clientNum++
			go sendData(conn, i, ch)
		}
	}
}

func main() {

	ch := make(chan int)
	connectServer(ch)

	go func() {
		for msg_num := range ch {
			msgNum += msg_num
		}
	}()

	time.Sleep(10 * time.Minute)
	fmt.Printf("rec msg %d \n", msgNum)
}
