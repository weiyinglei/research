package main

import (
    "fmt"
    "io"
    "math/rand"
    "net"
    "os"
    "time"
)

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}

func main() {
    ln, err := net.Listen("tcp", ":8080")
    checkError(err)

    for {
        conn, err := ln.Accept()
        if err != nil {
            fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
            continue
        }
        go run(conn)

    }
}

func run(conn net.Conn) {
    var all_len int = 0
    buffer := make([]byte, 20480)
    rand.Seed(time.Now().Unix())
    filename := string(rand.Intn(100))
    writeFile, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
    if err != nil {
        fmt.Println("create file error:", filename, err)
        return
    }
    defer writeFile.Close()
    for {

        lenght, err := conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Fprintf(os.Stderr, "read buff error: %s", err.Error())
            return
        }
        fmt.Println("receive data lenght:", lenght)

        _, err = writeFile.Write(buffer[:lenght])
        if err != nil {
            fmt.Println("write file error", err)
            return
        }

        all_len += lenght

    }
    fmt.Println("write file done", all_len/1024/1024)
    return
}
