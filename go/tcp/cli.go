package main

import (
    //"bytes"
    "fmt"
    //"io"
    "net"
    "os"
    //"time"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func main() {
    ch := make(chan int)
    for i := 0; i < 1000; i++ {
        go send(ch)
    }
    for w := range ch {
        fmt.Println(w)
    }
}

func send(ch chan int) {

    conn, err := net.Dial("tcp", "127.0.0.1:9999")
    checkError(err)
    //fileHandle, err := os.Open("app.log")
    //if err != nil {
    //    fmt.Println("open file ERROR", err)
    //    return
    //}
    //defer fileHandle.Close()
    //io.Copy(conn, fileHandle)
    conn.Write([]byte{'H'})
    //conn.Close()
    ch <- 1
}
