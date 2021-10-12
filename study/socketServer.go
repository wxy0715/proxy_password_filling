package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

func main() {
	if listener, err := net.Listen("tcp","127.0.0.1:9999"); err != nil {
		log.Fatal(err)
	} else {
		for true {
			conn, err := listener.Accept()
			if err != nil {
				log.Print(err) // 例如，连接终止
				continue
			}
			go handleConn(conn)
		}
	}
}

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	for input.Scan() {
		go echo(c, input.Text(), 1*time.Second)
	}
	// 注意：忽略 input.Err() 中可能的错误
	c.Close()
}

func handleConn0(c net.Conn) {
	io.Copy(c, c)
	c.Close()
}