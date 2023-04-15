package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	for {
		reqData := make([]byte, 1024)
		reqSize, err := conn.Read(reqData)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("err reading from client: ", err.Error())
				os.Exit(1)
			}
		}
		var resMessage string
		reqStrings := strings.TrimSpace(string(reqData[:reqSize]))

		if reqStrings == "PING" {
			resMessage = "+PONG\r\n"
		} else if len(reqStrings) >= 4 && reqStrings[0:4] == "PING" {
			resMessage = "+" + reqStrings[5:] + "\r\n"
		} else {
			resMessage = "+PONG\r\n"
		}
		resData := []byte(resMessage)
		_, err = conn.Write(resData)
		if err != nil {
			fmt.Println("Error sending response: ", err.Error())
		}
	}
}
