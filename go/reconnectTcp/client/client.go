package main

import (
	"net"
	"fmt"
	"bufio"
	"time"
)

var (
	hostAndPort = "127.0.0.1:12345"

	DefaultConnNum = 5000
)

func doRequest(conn net.Conn,connChan chan int,msgChan chan int64) {
	defer conn.Close()
	//ConnNum ++
	connChan <- 1
	for {
		fmt.Fprintf(conn,"HELLO WORLD\n")
		//MsgNum ++
		msgChan <- 1
		_,err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("recv data error")
			break
		}else{
			//fmt.Print("recv msg : ",msg)
		}
		time.Sleep(10 * time.Second)
	}
	//ConnNum --
	connChan <- -1
}

func newConnection(connChan chan int,msgChan chan int64){
	conn,err := net.Dial("tcp", hostAndPort)
	//fmt.Print("connect (", hostAndPort)
	if err != nil {
		fmt.Println("connect (", hostAndPort,") fail")
	}else{
		//fmt.Println(") ok",conn.LocalAddr(),"->",conn.RemoteAddr())
		go doRequest(conn,connChan,msgChan)
	}
}

func main() {
	var connNum int = 0
	var msgNum int64 = 0

	connChan := make(chan int)
	msgChan := make(chan int64)

	//print connNum coroutine
	go func() {
		for w := range connChan {
			//fmt.Println(w)
			connNum += w
			fmt.Println("connNum:",connNum)
		}
	}()
	//print msgNum coroutine
	go func() {
		for w := range msgChan {
			//fmt.Println(w)
			msgNum += w
			if msgNum % 100000 == 0 {
				fmt.Println("connNum:",connNum,"msgNum:",msgNum)
			}
		}
	}()
	//print connMsg and msgNum coroutine
	/*go func() {
		for {
			fmt.Println("connNum:", connNum,";msgNum:", msgNum)
			time.Sleep(3 * time.Second)
		}
	}()*/

	for {
		newConnection(connChan,msgChan)
		if connNum > (DefaultConnNum - 10) {
			time.Sleep(1 * time.Millisecond)
		}
		if connNum >= DefaultConnNum {
			break
		}
	}

	for {
		if connNum < DefaultConnNum {
			fmt.Println("ConnNum:DefaultConnNum:", connNum,":",DefaultConnNum)
			for i:=0;i<(DefaultConnNum- connNum);i++{
				newConnection(connChan,msgChan)
			}
		}
		time.Sleep(10 * time.Second)
	}
}