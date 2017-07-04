package main

import (
	"net"
	"time"
	"log"
	"os"
	"bytes"
	//"math/rand"
	"bufio"
)

var (
	hostAndPort = "10.0.30.231:8000"
	DefaultConnNum = 60000

	logger *log.Logger
)

func iovInfo() []byte {
	var bodyBuf bytes.Buffer
	bodyBuf.Write([]byte("[120.1767182, 30.1910232]"))
	bodyBuf.Write([]byte{0x0})
	bodyBuf.Write([]byte("LGWFF4V5840000001"))
	bodyBuf.Write([]byte{0x0})
	bodyBuf.Write([]byte("40"))
	bodyBuf.Write([]byte{0x0})
	bodyBuf.Write([]byte("2000"))
	bodyBuf.Write([]byte{0x0})
	bodyBuf.Write([]byte("100"))
	bodyBuf.Write([]byte{0x0})
	bodyBuf.Write([]byte("330102"))
	bodyBuf.Write([]byte{0x0})

	leng := len(bodyBuf.String())

	m := 0
	hex := make([]int, 0)
	length := 0;
	for{
		m = leng / 256
		leng = leng % 256

		if(m == 0){
			hex = append(hex, leng)
			length++
			break
		}

		hex = append(hex, m)
		length++;
	}

	var buf bytes.Buffer
	var b = []byte{1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 235}
	for i:=0; i<4-length;i++  {
		b = append(b,0)
	}

	for i:=0;i<len(hex) ;i++  {
		bTemp := byte(hex[i])
		b =append(b,bTemp)
	}

	buf.Write(b)
	buf.Write(bodyBuf.Bytes())
	return buf.Bytes()
}

func doRequest(conn net.Conn,connChan chan int,msgChan chan int64) {
	defer conn.Close()
	connChan <- 1
	for {
		b := iovInfo()

		//log.Println("buf:", string(b))
		conn.Write(b)
		msgChan <- 1
		_,err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Println("recv data error")
			break
		}else{
			//logger.Print("recv msg : ",msg)
		}
		//rand.Seed(int64(time.Now().Nanosecond()))
		//time.Sleep(time.Duration(rand.Intn(60000)) * time.Millisecond)
		time.Sleep(60 * time.Second)
	}
	connChan <- -1
}

func newConnection(connChan chan int,msgChan chan int64){
	conn,err := net.Dial("tcp", hostAndPort)
	if err != nil {
		logger.Println("connect (", hostAndPort,") fail")
		time.Sleep(1 * time.Second)
	}else{
		//logger.Println("connect (", hostAndPort,") ok",conn.LocalAddr(),"->",conn.RemoteAddr())
		go doRequest(conn,connChan,msgChan)
	}
}

func main() {
	//init logger
	file, err := os.Create("client.log")
	if err != nil {
		log.Fatalln("fail to create server.log file!")
	}
	logger = log.New(file, "", log.LstdFlags|log.Llongfile)

	//connNum,msgNum,connChan,msgChan
	var connNum int = 0
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
			logger.Println("ConnNum:DefaultConnNum:", connNum,":",DefaultConnNum)
			for i:=0;i<(DefaultConnNum- connNum);i++{
				newConnection(connChan,msgChan)
			}
		}
		time.Sleep(10 * time.Second)
	}
}
