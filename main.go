package main

import (
	"fmt"
	"os"
	"log"
	"strings"
)

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal(err)
	}

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
			fmt.Printf("read: %s\n", currentLine)
			newLineData = data[len(parts[0])+1:]
			currentLine = ""
		}
		currentLine += string(newLineData)
	}
}
