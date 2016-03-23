package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/loadimpact/speedboat/client"
	"github.com/loadimpact/speedboat/master"
	"time"
)

func init() {
	registerCommand(cli.Command{
		Name:   "ping",
		Usage:  "Test command, will be removed",
		Action: actionPing,
	})
	registerHandler(handlePing)
}

func handlePing(m *master.Master, msg master.Message, out chan master.Message) bool {
	switch msg.Type {
	case "ping.ping":
		out <- master.Message{
			Type: "ping.pong",
			Body: msg.Body,
		}
		return true
	default:
		return false
	}
}

func actionPing(c *cli.Context) {
	client, err := client.New("tcp://127.0.0.1:9595", "tcp://127.0.0.1:9596")
	if err != nil {
		log.WithError(err).Fatal("Couldn't create a client")
	}
	in, out, errors := client.Connector.Run()

	// Send a bunch of noise to filter through
	out <- master.Message{Type: "ping.noise"}
	out <- master.Message{Type: "ping.noise"}

	// Send a ping message, server should reply with a pong
	out <- master.Message{
		Type: "ping.ping",
		Body: time.Now().Format("15:04:05 2006-01-02 MST"),
	}

readLoop:
	for {
		select {
		case msg := <-in:
			switch msg.Type {
			case "ping.pong":
				log.WithFields(log.Fields{
					"body": msg.Body,
				}).Info("Pong!")
				break readLoop
			}
		case err := <-errors:
			log.WithError(err).Error("Ping failed")
		}
	}
}