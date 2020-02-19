package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

type params struct {
	Concurrent int
	LogDir     string
	NetAddr    net.IP
	NetMask    net.IPMask
	Password   string
	Ports      []uint
	Username   string
}

func (p *params) Dump() string {
	return fmt.Sprintf("Concurrent=%d NetAddr=%s NetMask=%s Password=%s Ports=%d Username=%s",
		p.Concurrent, p.NetAddr, p.NetMask, p.Password, p.Ports, p.Username)
}

func NewParams() (*params, error) {
	var netAddr cmdtoolbox.IPParam
	var netMask cmdtoolbox.IPMaskParam
	var ports portsParam

	flag.Var(&netAddr, "addr", "IP address of the network")
	concurrent := flag.Int("concurrent", 1, "sets the number of concurrent workers")
	logDir := flag.String("logdir", "", "path to the logs directory")
	flag.Var(&netMask, "mask", "IP address of the network mask")
	password := flag.String("password", "", "password for the DVR")
	flag.Var(&ports, "port", "port number")
	username := flag.String("username", "admin", "username for the DVR")

	flag.Parse()

	if logDir == nil || *logDir == "" {
		return nil, fmt.Errorf("specify logs directory")
	}

	if netAddr == nil {
		return nil, fmt.Errorf("specify IP address of the network")
	}

	if netMask == nil {
		return nil, fmt.Errorf("specify IP address of the network mask")
	}

	if len(ports) == 0 {
		return nil, fmt.Errorf("specify ports to scan")
	}

	return &params{
		Concurrent: *concurrent,
		LogDir:     *logDir,
		NetAddr:    net.IP(netAddr),
		NetMask:    net.IPMask(netMask),
		Password:   *password,
		Ports:      ports,
		Username:   *username,
	}, nil
}

type portsParam []uint

func (param *portsParam) String() string {
	return "port parameters"
}

func (param *portsParam) Set(value string) error {
	v, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return err
	}

	*param = append(*param, uint(v))

	return nil
}
