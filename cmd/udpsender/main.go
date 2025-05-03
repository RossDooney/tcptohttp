package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	udpAddr, err := net.ResolveUDPAddr("udp", "localhost:42069")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)

	defer conn.Close()

	if err != nil {
		fmt.Println("Errors :", err)
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {

		fmt.Print(">")
		text, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Errors : ", nil)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(text))

		if err != nil {
			fmt.Println("Errors :", err)
			os.Exit(1)
		}
	}
	// for {
	// 	conn, err := listener.Accept()

	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}

	// 	linesChan := getLinesChannel(conn)

	// 	for line := range linesChan {
	// 		fmt.Println("read:", line)
	// 	}
	// }
}

// func getLinesChannel(f io.ReadCloser) <-chan string {
// 	lines := make(chan string)
// 	go func() {
// 		defer f.Close()
// 		defer close(lines)
// 		currentLineContents := ""
// 		for {
// 			b := make([]byte, 8)
// 			n, err := f.Read(b)
// 			if err != nil {
// 				if currentLineContents != "" {
// 					lines <- currentLineContents
// 				}
// 				if errors.Is(err, io.EOF) {
// 					break
// 				}
// 				fmt.Printf("error: %s\n", err.Error())
// 				return
// 			}
// 			str := string(b[:n])
// 			parts := strings.Split(str, "\n")
// 			for i := 0; i < len(parts)-1; i++ {
// 				lines <- fmt.Sprintf("%s%s", currentLineContents, parts[i])
// 				currentLineContents = ""
// 			}
// 			currentLineContents += parts[len(parts)-1]
// 		}
// 	}()
// 	return lines
// }
