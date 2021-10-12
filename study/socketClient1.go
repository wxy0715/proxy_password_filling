package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go mustCopy1(os.Stdout, conn)
	mustCopy1(conn, os.Stdin)
}

func mustCopy1(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}