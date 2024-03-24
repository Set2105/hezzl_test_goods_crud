package nats

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
)

type NatsSettings struct {
	Host     string `xml:"natsHost"`
	Port     string `xml:"natsPort"`
	User     string `xml:"natsUser"`
	Password string `xml:"natsPassword"`
}

func (settings *NatsSettings) Valid() error {
	if settings.Host == "" {
		settings.Host = "nats"
	}
	if settings.Port == "" {
		settings.Port = "4222"
	}
	return nil
}

type Nats struct {
	Conn *nats.Conn
}

func InitNats(s *NatsSettings) (*Nats, error) {
	if err := s.Valid(); err != nil {
		return nil, fmt.Errorf("InitNats: %s", err)
	}
	nc, err := nats.Connect("nats://" + s.Host + ":" + s.Port)
	if err != nil {
		return nil, fmt.Errorf("InitNats: %s", err)
	}
	return &Nats{Conn: nc}, nil
}

func (n *Nats) Publish(subj string, data []byte) (err error) {
	err = n.Conn.Publish(subj, data)
	if err != nil {
		return fmt.Errorf("Nats.SendData.Publish: %s", err)
	}
	return nil
}
