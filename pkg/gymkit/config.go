package gymkit

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

type Rabbithost struct {
	RabbitHost, RabbitPort, RabbitUser, RabbitPassword string
}

// GetRabbitConfig builds an AMQP connection URL from environment variables.
// Supports RABBIT_HOST, RABBIT_HOST1, RABBIT_HOST2, … for round-robin selection.
func GetRabbitConfig() (connURL string, err error) {
	hosts := map[int]Rabbithost{
		0: {
			RabbitHost:     os.Getenv("RABBIT_HOST"),
			RabbitPort:     os.Getenv("RABBIT_PORT"),
			RabbitUser:     os.Getenv("RABBIT_USER"),
			RabbitPassword: os.Getenv("RABBIT_PASSWORD"),
		},
	}

	if hosts[0].RabbitHost == "" {
		return "", errors.New("RABBIT_HOST environment variable is not set")
	}
	if hosts[0].RabbitUser == "" {
		return "", errors.New("RABBIT_USER environment variable is not set")
	}

	for i := 1; ; i++ {
		h := os.Getenv("RABBIT_HOST" + strconv.Itoa(i))
		if h == "" {
			break
		}
		hosts[i] = Rabbithost{
			RabbitHost:     h,
			RabbitPort:     os.Getenv("RABBIT_PORT" + strconv.Itoa(i)),
			RabbitUser:     os.Getenv("RABBIT_USER" + strconv.Itoa(i)),
			RabbitPassword: os.Getenv("RABBIT_PASSWORD" + strconv.Itoa(i)),
		}
	}

	h := hosts[rand.Intn(len(hosts))]
	connURL = fmt.Sprintf("amqp://%s:%s@%s:%s/", h.RabbitUser, h.RabbitPassword, h.RabbitHost, h.RabbitPort)
	return connURL, nil
}
