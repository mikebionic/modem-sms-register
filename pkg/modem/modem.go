package modem

import (
	"context"
	"fmt"
	"time"

	"github.com/mike-bionic/modem-sms-register/pkg/config"
	log "github.com/sirupsen/logrus"
	"github.com/warthog618/modem/at"
	"github.com/warthog618/modem/gsm"
	"github.com/warthog618/modem/serial"
)

type Modem struct {
	gsm *gsm.GSM
}

func New(cfg *config.Config) (*Modem, error) {
	m, err := serial.New(serial.WithPort(cfg.SerialPort), serial.WithBaud(cfg.Baud))
	if err != nil {
		return nil, fmt.Errorf("failed to open serial port: %w", err)
	}

	g := gsm.New(at.New(m, at.WithTimeout(400*time.Millisecond)))
	if err := g.Init(); err != nil {
		m.Close()
		return nil, fmt.Errorf("failed to initialize GSM modem: %w", err)
	}

	return &Modem{gsm: g}, nil
}

func (m *Modem) Close() {
	if m.gsm != nil {
		m.gsm.StopMessageRx()
	}
}

func (m *Modem) StartMessageReceiver(ctx context.Context, url string) {
	m.gsm.StartMessageRx(
		func(msg gsm.Message) {
			log.WithFields(log.Fields{
				"from":    msg.Number,
				"message": msg.Message,
			}).Info("Received SMS message")

			if err := sendToEndpoint(url, msg.Number, msg.Message); err != nil {
				log.WithError(err).Error("Failed to send SMS to endpoint")
			}
		},
		func(err error) {
			log.WithError(err).Error("SMS receiver error")
		},
	)

	// polling signal quality
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if info, err := m.gsm.Command("+CSQ"); err != nil {
					log.WithError(err).Warn("Failed to read signal quality")
				} else {
					log.WithField("signal_quality", info).Debug("Signal quality update")
				}
			case <-m.gsm.Closed():
				return
			}
		}
	}()
}

func sendToEndpoint(url, number, message string) error {
	return nil
}
