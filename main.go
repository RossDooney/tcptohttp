package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {

	data, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	txtConent := ""
	for {
		b := make([]byte, 8)
		n, err := data.Read(b)

		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println(err)
			break
		}

		if err != nil {
			break
		}
		str := string(b[:n])
		txtConent += str
	}

	stuff := strings.Split(txtConent, "\n")
	for i := 0; i < len(stuff); i++ {
		if stuff[i] == "" {
			break
		}
		fmt.Printf("read: %v \n", stuff[i])
	}

}
