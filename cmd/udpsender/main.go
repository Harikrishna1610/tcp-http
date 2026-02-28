package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	conn, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	p, err1 := net.DialUDP("udp", nil, conn)
	if err1 != nil {
		log.Fatal(err1)
	}
	defer p.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		_, err = p.Write([]byte(line))

		if err != nil {
			log.Fatal(err)
		}

	}
}
