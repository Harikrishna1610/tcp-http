package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println(err)
	}
	for s := range getLinesChannel(f) {
		fmt.Printf("read: %s\n", s)
	}

}
func getLinesChannel(f io.ReadCloser) <-chan string {
	s := make(chan string, 1)
	go func() {
		defer f.Close()
		defer close(s)
		strs := ""
		for {
			var str []byte = make([]byte, 8)
			n, err := f.Read(str)
			if err != nil {
				break
			}
			str = str[:n]
			if i := bytes.IndexByte(str, '\n'); i != -1 {
				strs += string(str[:i])
				str = str[i+1:]
				s <- strs
				//fmt.Printf("read: %s\n", strs)
				strs = ""
			}
			strs += string(str)

		}
		if len(strs) != 0 {
			s <- strs
		}
	}()
	return s
}
