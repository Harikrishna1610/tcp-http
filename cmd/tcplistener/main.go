package main

import (
	"fmt"
	"http_frm_udp/internal/request"
	"net"
)

func main() {
	// this is a tcp listener on port
	f, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Unable to Open \n")
	}
	defer f.Close()
	for true {
		conn, err := f.Accept() // upon successful listen this writes to CLI
		//fmt.Println("listening on ", f.Addr())
		if err != nil {
			fmt.Println("Unable to get connection")
		}

		res, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("unable to read from parser")
		}
		fmt.Println("Request Line:")
		fmt.Println("Method:", res.RequestLine.Method)
		fmt.Println("Target:", res.RequestLine.RequestTarget)
		fmt.Println("Version:", res.RequestLine.HttpVersion)

		// for line := range getLinesChannel(conn) { //we can directly send conn to our old method
		// 	fmt.Println(line)
		// }

	}

}

// func getLinesChannel(f io.ReadCloser) <-chan string { //this return type is called channel of strings
// 	var s chan string = make(chan string, 1) // create a channel of string with length 1
// 	go func() {                              // beggining of go routines
// 		defer f.Close() //using defer close so that the connections can be closed at the end
// 		defer close(s)
// 		var buildString = ""

// 		for {
// 			var b []byte = make([]byte, 8)
// 			n, err := f.Read(b)
// 			if err != nil {
// 				break
// 			}
// 			b = b[:n] // we are doing this because this is an 8 size array so if the string has 2 chars the rest 6 items in array are older line's which make it cumbersome
// 			if i := (bytes.IndexByte(b, '\n')); i != -1 {
// 				buildString += string(b[:i])
// 				b = b[i+1:]
// 				s <- buildString //writing back to channel is with this
// 				buildString = ""
// 			}
// 			buildString += string(b)
// 		}
// 		if len(buildString) != 0 {
// 			s <- buildString
// 		}

// 	}()
// 	//var parts []string
// 	return s

// }
