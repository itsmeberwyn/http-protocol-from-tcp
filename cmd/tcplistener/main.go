package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

// i started working on the file accessing method as it
// is closely related to creating a similar approach on
// handling request
func acceptDataFromFile() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Connection has been Accepted!")
		for line := range getLinesChannel(conn) {
			fmt.Printf("%s\n", line)
		}
		log.Println("Connection has been Closed!")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		currentLine := ""
		for {
			data := make([]byte, 8)
			count, err := f.Read(data)
			if err != nil {
				break
			}

			parts := strings.Split(string(data[:count]), "\n")
			newLineData := data

			if len(parts) > 1 {
				currentLine += parts[0]
				out <- currentLine
				newLineData = data[len(parts[0])+1:]
				currentLine = ""
			}
			currentLine += string(newLineData)
		}
	}()

	return out
}
