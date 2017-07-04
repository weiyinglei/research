package main

import (
	"net"
	"os"
	"strings"
	"log"
	"github.com/reddec/go-queue"
	"fmt"
)

var (
	port = ":8090"

	logger *log.Logger
	q *queue.BlockingQueue
)
func checkError(err error) {
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}
}
func handleClient(conn net.Conn,connChan chan int,msgChan chan int64) {
	//conn.SetReadDeadline(time.Now().Add(3 * time.Minute)) error
	request := make([]byte,1024)
	defer conn.Close()

	connChan <- 1
	for {
		recv_len,err := conn.Read(request)
		if err != nil {
			logger.Println(err)
			break
		}
		if recv_len == 0 {
			break
		}
		recvData := strings.TrimSpace(string(request[:recv_len]))
		msgChan <- 1
		q.Put(recvData)

		//logger.Println("recv_len : ",recv_len,",recv_data : " + recvData)
		conn.Write([]byte("O\n"))
		request = make([]byte,1024)
	}
	connChan <- -1
}
func main() {
	//init logger
	file, err := os.Create("server.log")
	if err != nil {
		log.Fatalln("fail to create server.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)

	//init queue
	q = queue.New()

	//start tcp server
	tcpAddr,err := net.ResolveTCPAddr("tcp4", port)
	checkError(err)
	listener,err := net.ListenTCP("tcp",tcpAddr)
	defer listener.Close()
	checkError(err)
	logger.Println("server is running at ", port)

	//connNum,msgNum,connChan,msgChan
	connNum := 0
	var msgNum int64 = 0

	connChan := make(chan int)
	msgChan := make(chan int64)
	//print connNum coroutine
	go func() {
		for w := range connChan {
			connNum += w
			logger.Println("connNum:",connNum)
		}
	}()
	//print connNum,msgNum coroutine
	go func() {
		for w := range msgChan {
			msgNum += w
			if msgNum % 100000 == 0 {
				logger.Println("connNum:",connNum,"msgNum:",msgNum)
			}
		}
	}()

	go func() {
		for {
			str,b := q.Pop();
			if b {
				fmt.Println(str.(string))
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
