package main

// Go provides a `flag` package supporting basic
// command-line flag parsing. We'll use this package to
// implement our example command-line program.
import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.bug.st/serial"
	"golang.org/x/term"
)

func main() {

	portListFlag := flag.Bool("l", false, "a bool")
	portSelFlag := flag.String("port", "COM3", "a string")
	portBaudFlag := flag.Int("baud", 115200, "an int")

	flag.Parse()

	if *portListFlag {
		ports, err := serial.GetPortsList()
		if err != nil {
			log.Fatal(err)
		}
		if len(ports) == 0 {
			log.Fatal("No serial ports found!")
		}
		for _, port := range ports {
			fmt.Printf("Found port: %v\n", port)
		}
	} else {
		mode := &serial.Mode{
			BaudRate: *portBaudFlag,
		}
		port, err := serial.Open(*portSelFlag, mode)
		port.SetReadTimeout(250)
		fmt.Printf("Opening port: %v at %v baud rate\n", *portSelFlag, *portBaudFlag)
		fmt.Println("Use <Ctrl> A to exit")
		go func() {
			for {
				// switch stdin into 'raw' mode
				oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
				if err != nil {
					fmt.Println(err)
					return
				}
				defer term.Restore(int(os.Stdin.Fd()), oldState)

				b := make([]byte, 1)
				_, err = os.Stdin.Read(b)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("%c", (b[0]))
				port.Write(b)
				if b[0] == 1 {
					fmt.Println("Goodbye")
					os.Exit(0)
				}
			}
		}()

		if err != nil {
			log.Fatal(err)
		}
		x := 0
		for {
			//Read Serail Port
			buf := make([]byte, 128)
			x, err = port.Read(buf)
			if err != nil {
				//lost serial port
				log.Fatal(err)
			}

			if x == -1 {
				fmt.Println("should not be here")
			}

			strbuf := string(buf[:])
			fmt.Print(strbuf)

		}
	}
}
