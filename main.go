package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/morgulbrut/color256"
	"github.com/tarm/serial"
)

var c serial.Config
var s *serial.Port

func main() {
	reader := bufio.NewReader(os.Stdin)

	flag.StringVar(&c.Name, "port", "/dev/ttyUSB0", "set port")
	flag.IntVar(&c.Baud, "baud", 115200, "set baudrate")
	durSt := flag.String("timeout", "3s", "set read timeout. format: 5s, 500ms")
	par := flag.String("parity", "N", "sets the parity. possible values: N, O, E, M, S")
	stopb := flag.Int("stopbits", 1, "sets the stopbits. possible values: 1, 15, 2")
	datas := flag.Int("datasize", 8, "sets the datasize")

	flag.Parse()

	c.ReadTimeout, _ = time.ParseDuration(*durSt)
	c.Parity = serial.Parity([]byte(*par)[0])
	c.StopBits = serial.StopBits(*stopb)
	c.Size = byte(*datas)

	printSettings()
	fmt.Println()
	color256.PrintHiCyan("/help:\t for help")

	s, err := serial.OpenPort(&c)
	if err != nil {
		color256.PrintHiRed(err.Error())
		os.Exit(1)
	}

	defer s.Close()

	go read(s)

	for {
		cmdString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err.Error())
			os.Exit(2)
		}
		runCommand(cmdString, s)
	}
}

func runCommand(commandStr string, s *serial.Port) {
	commandStr = strings.TrimSuffix(commandStr, "\n")
	arrCommandStr := strings.Fields(commandStr)
	if len(arrCommandStr) > 0 {
		switch arrCommandStr[0] {
		case "/exit":
			os.Exit(0)
		case "/quit":
			os.Exit(0)
		case "/settings":
			printSettings()
		case "/help":
			printSettings()
			printHelp()
		default:
			cmd := []byte(commandStr + "\n\r")
			_, err := s.Write(cmd)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}
}

func printHelp() {
	color256.PrintHiGreen("======== COMMANDS =========")
	fmt.Printf("/exit:\t\t exits the programm\n")
	fmt.Printf("/quit:\t\t exits the programm\n")
	fmt.Printf("/settings:\t shows the settings\n")
	fmt.Printf("/help:\t\t shows this help\n")
}

func printSettings() {
	color256.PrintHiGreen("======== SETTINGS =========")
	fmt.Printf("Port:\t\t%s\n", c.Name)
	fmt.Printf("Baud:\t\t%d\n", c.Baud)
	fmt.Printf("Parity:\t\t%q\n", c.Parity)
	fmt.Printf("Stopbits:\t%d\n", c.StopBits)
	fmt.Printf("Datasize:\t%d\n", c.Size)
	fmt.Printf("Timeout:\t%q\n", c.ReadTimeout)
}

func read(s *serial.Port) {
	for true {
		buf := make([]byte, 512)
		n, err := s.Read(buf)
		if err != nil {
			if runtime.GOOS == "windows" {
				color256.PrintHiRed("[read]: " + err.Error())
				os.Exit(2)
			}
		}
		if n != 0 {
			fmt.Print(string(buf[:n]))
		}
	}
}
