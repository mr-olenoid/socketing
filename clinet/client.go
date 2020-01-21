package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const bufferSize = 1024

func readIO(conn net.Conn) {
	fileSize := 117338112
	fmt.Println("Reciving file: ", "file.msi")
	buf := make([]byte, bufferSize)
	f, err := os.OpenFile("file.msi", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer f.Close()
	if err != nil {
		log.Println("Error opening file, ", err)
		return
	}
	for fileSize > 0 {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading stream: ", err)
				return
			}
		}
		_, ferr := f.Write(buf[:n])
		if ferr != nil {
			log.Println("Error writung to file, ", ferr)
			return
		}
		fileSize -= n
	}
	fmt.Println("Done!!")
}

func writeIO(conn net.Conn) {
	//buf := make([]byte, 0, 256)
	//buf := []byte("Hello world")
	f, err := os.Open("toSend.msi")
	if err != nil {
		log.Println("Error opening file, ", err)
		return
	}
	/*
		fi, err := f.Stat()
		if err != nil {
			// Could not obtain stat, handle error
		}
		fmt.Printf("The file is %d bytes long", fi.Size())
	*/
	defer f.Close()
	buf := make([]byte, bufferSize)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println("Error reading file ", err)
				return
			}
		}
		_, cerr := conn.Write(buf[:n])
		if cerr != nil {
			log.Println("Connection error, ", cerr)
			return
		}
		if n < bufferSize {
			fmt.Println("Done!!")
			break
		}
	}
}

func main() {
	var conn *net.TCPConn

	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8000")
	if err != nil {
		println("ResolveTCPAddr failed:", err.Error())
		os.Exit(1)
	}
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
	}
	conn.SetKeepAlive(true)
	conn.SetKeepAlivePeriod(10)

	go writeIO(conn)
	go readIO(conn)
	for {
		time.Sleep(60)
	}
}
