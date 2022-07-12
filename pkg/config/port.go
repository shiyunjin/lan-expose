package config

import (
	"errors"
	"strconv"
)

type Port struct {
	port string
}

func (p *Port) Set(port string) error {
	// valid port
	porti, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	if 0 < porti && porti <= 65535 {
		p.port = port
		return nil
	} else {
		return errors.New("port is valid, range 0 < port <= 65535")
	}
}

func (p *Port) SafeSet(port string) error {
	// valid port with safe rule
	porti, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	if 80 < porti && porti <= 1024 && porti != 443 && porti != 8443 {
		p.port = port
		return nil
	} else {
		return errors.New("port valid range is 80 < port <= 1024 and not 443,8443")
	}
}

func (p Port) Get() string {
	return p.port
}
