package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
	"github.com/warthog618/modem/trace"
)

var version = "undefined"

func main() {

	serial_port, url_address, err := get_port_and_url_from_config()
	if err != nil {
		fmt.Println(err)
	}

	dev := flag.String("d", serial_port, "path to modem device")
	baud := flag.Int("b", 115200, "baud rate")
	period := flag.Duration("p", 100*time.Minute, "period to wait")
	timeout := flag.Duration("t", 400*time.Millisecond, "command timeout period")
	verbose := flag.Bool("v", false, "log modem interactions")
	hex := flag.Bool("x", false, "hex dump modem responses")
	url := flag.String("u", url_address, "sms register endpoint url")
	vsn := flag.Bool("version", false, "report version and exit")
	flag.Parse()

	if *vsn {
		fmt.Printf("%s %s\n", os.Args[0], version)
		os.Exit(0)
	}
	m, err := serial.New(serial.WithPort(*dev), serial.WithBaud(*baud))
	if err != nil {
		log.Println(err)
		return
	}
	defer m.Close()
	var mio io.ReadWriter = m
	if *hex {
		mio = trace.New(m, trace.WithReadFormat("r: %v"))
	} else if *verbose {
		mio = trace.New(m)
	}
	g := gsm.New(at.New(mio, at.WithTimeout(*timeout)))
	err = g.Init()
	if err != nil {
		log.Println(err)
		return
	}

	go pollSignalQuality(g, timeout)

	err = g.StartMessageRx(
		func(msg gsm.Message) {
			log.Printf("%s: %s\n", msg.Number, msg.Message)
			if reqerr := make_request(*url, msg.Number, msg.Message); reqerr != nil {
				fmt.Println(reqerr)
			}
		},
		func(err error) {
			log.Printf("err: %v\n", err)
		})

	if err != nil {
		log.Println(err)
		return
	}
	defer g.StopMessageRx()

	for {
		select {
		case <-time.After(*period):
			log.Println("exiting...")
			// return
		case <-g.Closed():
			log.Fatal("modem closed, exiting...")
		}
	}
}

// pollSignalQuality polls the modem to read signal quality every minute.
//
// This is run in parallel to SMS reception to demonstrate separate goroutines
// interacting with the modem.
func pollSignalQuality(g *gsm.GSM, timeout *time.Duration) {
	for {
		select {
		case <-time.After(time.Minute):
			i, err := g.Command("+CSQ")
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("Signal quality: %v\n", i)
			}
		case <-g.Closed():
			return
		}
	}
}
