package main

// golang实现带有心跳检测的tcp长连接
// server
import (
	"fmt"
	"net"
	"time"
)

// message struct:
// c#d

var (
	Req_REGISTER byte = 1 // 1 --- c register cid
	Res_REGISTER byte = 2 // 2 --- s response

	Req_HEARTBEAT byte = 3 // 3 --- s send heartbeat req
	Res_HEARTBEAT byte = 4 // 4 --- c send heartbeat res

	Req byte = 5 // 5 --- cs send data
	Res byte = 6 // 6 --- cs send ack
)

type CS struct {
	Rch chan []byte
	Wch chan []byte
	Dch chan bool
	u   string
}

func NewCs(uid string) *CS {
	return &CS{Rch: make(chan []byte), Wch: make(chan []byte), u: uid}
}

var CMap map[string]*CS

func main() {
	CMap = make(map[string]*CS)
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("127.0.0.1"), 6666, ""})
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}
	fmt.Println("已初始化连接，等待客户端连接...")
	go PushGRT()
	Server(listen)
	select {}
}

func PushGRT() {
	for {
		time.Sleep(15 * time.Second)
		for k, v := range CMap {
			fmt.Println("push msg to user:" + k)
			v.Wch <- []byte{Req, '#', 'p', 'u', 's', 'h', '!'}
		}
	}
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		// handler goroutine
		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	defer conn.Close()
	data := make([]byte, 128)
	var uid string
	var C *CS
	for {
		conn.Read(data)
		fmt.Println("客户端发来数据:", string(data))
		if data[0] == Req_REGISTER { // register
			conn.Write([]byte{Res_REGISTER, '#', 'o', 'k'})
			uid = string(data[2:])
			C = NewCs(uid)
			CMap[uid] = C
			//			fmt.Println("register client")
			//			fmt.Println(uid)
			break
		} else {
			conn.Write([]byte{Res_REGISTER, '#', 'e', 'r'})
		}
	}
	//	WHandler
	go WHandler(conn, C)

	//	RHandler
	go RHandler(conn, C)

	//	Worker
	go Work(C)
	select {
	case <-C.Dch:
		fmt.Println("close handler goroutine")
	}
}

// 正常写数据
// 定时检测 conn die => goroutine die
func WHandler(conn net.Conn, C *CS) {
	// 读取业务Work 写入Wch的数据
	ticker := time.NewTicker(20 * time.Second)
	for {
		select {
		case d := <-C.Wch:
			conn.Write(d)
		case <-ticker.C:
			if _, ok := CMap[C.u]; !ok {
				fmt.Println("conn die, close WHandler")
				return
			}
		}
	}
}

// 读客户端数据 + 心跳检测
func RHandler(conn net.Conn, C *CS) {
	// 心跳ack
	// 业务数据 写入Wch

	for {
		data := make([]byte, 128)
		// setReadTimeout
		err := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if err != nil {
			fmt.Println(err)
		}
		if _, derr := conn.Read(data); derr == nil {
			// 可能是来自客户端的消息确认
			//           	     数据消息
			fmt.Println(data)
			if data[0] == Res {
				fmt.Println("recv client data ack")
			} else if data[0] == Req {
				fmt.Println("recv client data")
				fmt.Println(data)
				conn.Write([]byte{Res, '#'})
				// C.Rch <- data
			}

			continue
		}

		conn.Write([]byte{Req_HEARTBEAT, '#'})
		fmt.Println("send ht packet")
		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		if _, herr := conn.Read(data); herr == nil {
			// fmt.Println(string(data))
			fmt.Println("resv ht packet ack")
		} else {
			delete(CMap, C.u)
			fmt.Println("delete user!")
			return
		}
	}
}

func Work(C *CS) {
	time.Sleep(5 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}

	time.Sleep(15 * time.Second)
	C.Wch <- []byte{Req, '#', 'h', 'e', 'l', 'l', 'o'}
	// 从读ch读信息
	/*	ticker := time.NewTicker(20 * time.Second)
		for {
			select {
			case d := <-C.Rch:
				C.Wch <- d
			case <-ticker.C:
				if _, ok := CMap[C.u]; !ok {
					return
				}
			}

		}
	*/ // 往写ch写信息
}
