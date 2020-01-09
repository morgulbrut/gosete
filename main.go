package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/morgulbrut/color256"
	"github.com/tarm/serial"
)

var c serial.Config

func main() {
	reader := bufio.NewReader(os.Stdin)

	c.Name = "/dev/ttyUSB0"
	c.Baud = 115200
	c.ReadTimeout = time.Second * 3
	c.Parity = 'N'
	c.StopBits = 1
	c.Size = 8

	printSettings()

	for {
		fmt.Print("> ")
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
		}
		runCommand(cmdString)
	}
}

func runCommand(commandStr string) {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) > 0 {
		switch arrCommandStr[0] {
		case "/exit":
			os.Exit(0)
		case "/quit":
			os.Exit(0)
		case "/port":
			c.Name = arrCommandStr[1]
		case "/baud":
			c.Baud, _ = strconv.Atoi(arrCommandStr[1])
		case "/timeout":
			c.ReadTimeout, _ = time.ParseDuration(arrCommandStr[1])
		case "/settings":
			printSettings()
		case "/help":
			printSettings()
		default:
			serCom(commandStr)
		}
	}
}

func printSettings() {
	fmt.Println("======== SETTINGS =========")
	fmt.Printf("Port:\t\t%s\n", c.Name)
	fmt.Printf("Baud:\t\t%d\n", c.Baud)
	fmt.Printf("Parity:\t\t%q\n", c.Parity)
	fmt.Printf("Stopbits:\t%d\n", c.StopBits)
	fmt.Printf("Datasize:\t%d\n", c.Size)
	fmt.Printf("Timeout:\t%q\n", c.ReadTimeout)
}

func serCom(msg string) {
	s, err := serial.OpenPort(&c)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer s.Close()

	n, err := s.Write([]byte(msg + "\n\r"))
	if err != nil {
		log.Fatal(err.Error())
	}
	for true {

		buf := make([]byte, 512)
		n, err = s.Read(buf)
		if err != nil {
			color256.Red(err.Error())
		}
		fmt.Print(string(buf[:n]))
		if n == 0 {
			break
		}
	}
}
