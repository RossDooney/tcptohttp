package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {

	data, err := os.Open("./messages.txt")
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	//br := bufio.NewReader(data)

	for {
		b := make([]byte, 8)
		_, err := data.Read(b)

		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Println(err)
			break
		}

		if err != nil {
			break
		}
		fmt.Printf("read: %v\n", string(b))

	}

}
