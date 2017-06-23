package main

import (
	"log"
	"os"
	"io/ioutil"
	"io"
	"fmt"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 致命错误
)

func init() {
	fmt.Println("init")
	file, err := os.OpenFile("file", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
	}

	Trace = log.New(ioutil.Discard, "TRACE: ", log.Ltime|log.Lshortfile)
	Info = log.New(os.Stdout, "Info: ", log.Ltime|log.Lshortfile)
	Warning = log.New(os.Stdout, "Warning: ", log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(file, os.Stderr),  "Error", log.Ltime|log.Lshortfile)
}

func main() {
	Trace.Println("I have something standard to say")
	Info.Println("Special Information")
	Warning.Println("There is something you need to know about")
	Error.Println("Something has failed")
}