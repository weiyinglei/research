package main

import (
	//"bytes"
	"fmt"
	//"io"
	"math/rand"
	//"net"
	"os"
	"time"
	"database/sql"
	"strconv"
	"container/list"
	"strings"
	"bytes"
	"net"
	//"runtime/debug"
	"net/url"
	_ "github.com/go-sql-driver/mysql"
	//_ "github.com/davecheney/gmx"

)
var (
	db *sql.DB
	lvin *list.List
	svin []string
	sCurrentVin []string
	lcoor *list.List
	scoor []string

	host string = "10.0.30.120:8090"
	//host string = "10.0.30.124:8090"
	//host string = "10.70.7.181:8090"

	hp = "tcp://10.0.30.120:8090"
	//hp = "tcp://10.0.30.215:8090"
	//hp = "tcp://10.70.7.215:8090"

	DefalutTimeout = 25 * time.Second
	MaxClient      = 10000
	clientNum      = 0
)
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
func open(ch chan int) net.Conn {
	//rand.Seed(int64(time.Now().Nanosecond()))
	//time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	var conn net.Conn = nil
	for {
		_, err := net.Dial("tcp", host)
		if err == nil {
			ch <- 1
			break
		}
	}

	defer func(ch chan int) {
		//fmt.Println("disconnected :" + ipStr)
		if conn != nil {
			err := conn.Close()
			if err != nil {
				fmt.Println(err)
				//fmt.Println(debug.Stack())
			}
			ch <- -1
		}
	}(ch)
	return conn
}
func write(conn net.Conn,buf bytes.Buffer){
	for {
		_,err := conn.Write([]byte(buf.String()))
		if err != nil {
			//fmt.Println("write:",conn)
			//fmt.Println(debug.Stack())
			//fmt.Println("write error!",err)
			time.Sleep(1000 * time.Millisecond)
			conn,err = sockConn(hp,DefalutTimeout)
			if err == nil{
				fmt.Println(conn.LocalAddr())
				continue
			}
		}
		if err == nil {
			//fmt.Println(conn.LocalAddr()," write: ",buf.String())
			break
		}
	}

}
func start(ch chan int, coor []string, sCurrentVin []string) {
	conn := open(ch)
	for {
		//r := rand.New(rand.NewSource(time.Now().UnixNano()))
		rand.Seed(int64(time.Now().Nanosecond()))
		time.Sleep(time.Duration(rand.Intn(10000) + 10000) * time.Millisecond)

		var bodyBuf bytes.Buffer
		bodyBuf.WriteString("{\"c\":\"[")
		bodyBuf.WriteString(coor[rand.Intn(len(coor))])
		bodyBuf.WriteString("]\",\"v\":\"")
		bodyBuf.WriteString(sCurrentVin[rand.Intn(len(sCurrentVin))])
		bodyBuf.WriteString("\"}")
		body := bodyBuf.String()
		leng := len(body)

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

		//fmt.Println(hex)

		var buf bytes.Buffer
		var b = []byte{1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7, 209}
		for i:=0; i<4-length;i++  {
			b = append(b,0)
		}

		for i:=0;i<len(hex) ;i++  {
			bTemp := byte(hex[i])
			b =append(b,bTemp)
		}

		//b = append(b,byte(leng))

		buf.Write(b)
		buf.WriteString(bodyBuf.String())

		write(conn,buf)
		/*if err != nil {
			time.Sleep(time.Duration(rand.Intn(1000) + 1000) * time.Millisecond)
			conn, err := net.Dial("tcp", host)
			if err == nil{
				ch <- 1
			}
			if err != nil {
				fmt.Println(err)
			}
			defer func(ch chan int) {
				//fmt.Println("disconnected :" + ipStr)
				conn.Close()
				ch <- -1
			}(ch)
		}*/
	}
}
func initSvin(){
	lvin = list.New()

	rows, err := db.Query("SELECT vin FROM vehicle")
	defer rows.Close()
	checkError(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	for rows.Next() {
		record := make(map[string]string)
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		lvin.PushBack(record)
	}

	for v := lvin.Front(); v != nil; v = v.Next() {
		for _,v := range v.Value.(map[string]string) {
			//fmt.Print(k,":",v)
			svin = append(svin,v)
		}
	}
}
func initSarea(){
	lcoor = list.New()

	rows, err := db.Query("SELECT coordinate FROM area WHERE leaf = '1'")
	defer rows.Close()
	checkError(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	for rows.Next() {
		record := make(map[string]string)
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		lcoor.PushBack(record)
	}

	for v := lcoor.Front(); v != nil; v = v.Next() {
		for _,v := range v.Value.(map[string]string) {
			//fmt.Print(k,":",v)
			strs := strings.Split(strings.Split(v, "[")[1],"]")
			scoor = append(scoor,strs[0])
		}
	}
}
func count() int {
	count := 0
	rows, err := db.Query("SELECT count(1) cnt FROM vehicle")
	defer rows.Close()
	checkError(err)

	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}
	for rows.Next() {
		record := make(map[string]string)
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		b,error := strconv.Atoi(record["cnt"])
		checkError(error)
		count = b
		return count
	}
	return count
}
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
	for {
		rand.Seed(int64(time.Now().Nanosecond()))
		time.Sleep(time.Duration(rand.Intn(60000)) * time.Millisecond)

		var bodyBuf bytes.Buffer
		bodyBuf.WriteString("{\"c\":\"[")
		bodyBuf.WriteString(scoor[rand.Intn(len(scoor))])
		bodyBuf.WriteString("]\",\"v\":\"")
		bodyBuf.WriteString(sCurrentVin[rand.Intn(len(sCurrentVin))])
		bodyBuf.WriteString("\"}")
		body := bodyBuf.String()
		leng := len(body)

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

		//fmt.Println(hex)

		var buf bytes.Buffer
		var b = []byte{1, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 235}
		for i:=0; i<4-length;i++  {
			b = append(b,0)
		}

		for i:=0;i<len(hex) ;i++  {
			bTemp := byte(hex[i])
			b =append(b,bTemp)
		}

		//b = append(b,byte(leng))

		buf.Write(b)
		buf.WriteString(bodyBuf.String())

		write(conn,buf)
	}
}
func connectServer(ch chan int) {
	for i := 1; i <= MaxClient; i++ {
		conn, err := sockConn(hp, DefalutTimeout)
		if err != nil {
			fmt.Printf("connect error:%s\n", err)
		} else {
			clientNum++
			go sendData(conn, i, ch)
		}
	}
}
func main() {
	//连接数据库
	db, _ = sql.Open("mysql", "root:123.com@tcp(10.0.30.120:3306)/iov?charset=utf8")
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.Ping()

	initSvin()

	initSarea()

	//vin列表
	vinCount := count()
	//s := make([]string,count,count)

	//随机vin
	// 根据时间设置随机数种子
	rand.Seed(int64(time.Now().Nanosecond()))
	// 获取指定范围内的随机数
	for i := 0; i < 100000; i ++ {
		index := rand.Intn(vinCount)
		sCurrentVin = append(sCurrentVin,svin[index])
	}

	count := 0
	ch := make(chan int)

	connectServer(ch)

	for w := range ch {
		//fmt.Println(w)
		count += w
		fmt.Println(count)
	}
}
