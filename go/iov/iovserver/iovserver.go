package main

import (
	"net"
	"os"
	"time"
	"github.com/mkideal/log"
)
const (
	MAX_CONN_NUM = 10
)
func EchoFunc(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 100)
	for {
		if conn != nil {
			_, err := conn.Read(buf)
			if err != nil {
				log.Error("Error reading:", conn.RemoteAddr(), err.Error())
				return
			}
			if err == nil {
				log.Error("listen ",conn.RemoteAddr(),string(buf))
			}
			//send reply
			_, err = conn.Write(buf)
			if err != nil {
				log.Error("Error send reply:", conn.RemoteAddr(), err.Error())
				return
			}
		}
	}
}
func initLog()  {
	// Init and defer Uninit
	defer log.Uninit(log.InitFileAndConsole("./logs/app.log",log.LvDEBUG))

	log.SetLevel(log.LvINFO)

	// 默认日志等级是 INFO, 可以按以下方式修改等级:
	//
	//	log.SetLevel(log.LvTRACE)
	// 	log.SetLevel(log.LvDEBUG)
	// 	log.SetLevel(log.LvINFO)
	// 	log.SetLevel(log.LvWARN)
	// 	log.SetLevel(log.LvERROR)
	// 	log.SetLevel(log.LvFATAL)
}
func init()  {
	initLog()
}
func main() {

	/*c, err := redis.Dial("tcp", "10.0.30.120:6379")
	if err != nil {
		fmt.Println(err)
		return
	}*/
	//defer c.Close()

	listener, err := net.Listen("tcp", "10.70.7.181:8090")
	if err != nil {
		log.Error("error listening:", err.Error())
		os.Exit(1)
	}
	defer listener.Close()

	log.Error("running ...\n")

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
			log.Error("connNum: %d\n", connNum)
		}
	}()

	for i := 0; i < MAX_CONN_NUM; i++ {
		go func() {
			for conn := range connChan {
				connChangeChan <- 1
				EchoFunc(conn)
				connChangeChan <- -1
			}
		}()
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Error("Error accept:", err.Error())
			return
		}
		connChan <- conn
	}
}