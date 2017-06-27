package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	//"time"
)

var (
	port = ":12345"
)

func checkError(err error) {
	if err != nil {
		//fmt.Println(err)
		os.Exit(1)
	}
}

func handleClient(conn net.Conn,connChan chan int,msgChan chan int64) {
	//conn.SetReadDeadline(time.Now().Add(3 * time.Minute))
	//conn.SetDeadline(time.Now().Add(1 * time.Minute))
	request := make([]byte,1024)
	defer conn.Close()

	//ConnNum ++
	connChan <- 1
	for {
		recv_len,err := conn.Read(request)
		if err != nil {
			fmt.Println(err)
			break
		}
		if recv_len == 0 {
			break
		}
		strings.TrimSpace(string(request[:recv_len]))
		//MsgNum ++
		msgChan <- 1
		//fmt.Println("recv_len : ",recv_len)
		//fmt.Println("recv_data : " + recvData)
		//daytime := time.Now().String()
		//conn.Write([]byte(daytime + "\n"))
		conn.Write([]byte("O\n"))
		request = make([]byte,1024)
	}
	//ConnNum --
	connChan <- -1
}

func main() {
	tcpAddr,err := net.ResolveTCPAddr("tcp4", port)
	checkError(err)
	listener,err := net.ListenTCP("tcp",tcpAddr)
	defer listener.Close()
	checkError(err)

	/*go func() {
		for {
			fmt.Println("ConnNum:",ConnNum,";MsgNum:",MsgNum)
			time.Sleep(3 * time.Second)
		}
	}()*/

	connNum := 0
	var msgNum int64 = 0

	connChan := make(chan int)
	msgChan := make(chan int64)
	go func() {
		for w := range connChan {
			//fmt.Println(w)
			connNum += w
			fmt.Println("connNum:",connNum)
		}
	}()
	go func() {
		for w := range msgChan {
			//fmt.Println(w)
			msgNum += w
			if msgNum % 100000 == 0 {
				fmt.Println("connNum:",connNum,"msgNum:",msgNum)
			}
		}
	}()

	for {
		conn,err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn, connChan, msgChan)
	}

}