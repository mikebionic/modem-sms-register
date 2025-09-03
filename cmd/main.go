package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/mike-bionic/modem-sms-register/pkg/config"
	"github.com/mike-bionic/modem-sms-register/pkg/modem"
	log "github.com/sirupsen/logrus"
)

var version = "1.0.0"

func main() {
	cfg, err := config.GetConfigData()
	if err != nil {
		log.Fatal(err)
	}

	dev := flag.Bool("dev", false, "ignore modem connection")
	verbose := flag.Bool("v", false, "enable verbose logs")
	flag.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.WithFields(log.Fields{
		"device": cfg.SerialPort,
		"baud":   cfg.Baud,
		"url":    cfg.GetURL(),
	}).Info("Starting SMS register service")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var m *modem.Modem
	if !*dev {
		m, err = modem.New(cfg)
		if err != nil {
			log.Fatal(err)
		}
		defer m.Close()
	} else {
		log.Info("Running without real modem")
	}

	if m != nil {
		go m.StartMessageReceiver(ctx, cfg.GetURL())
	}

	log.Info("Service started, waiting for messages...")

	select {
	case <-sigCh:
		log.Info("Shutdown signal received, exiting...")
	case <-ctx.Done():
		log.Info("Context canceled, exiting...")
	}

	log.Info("Service stopped")
}
