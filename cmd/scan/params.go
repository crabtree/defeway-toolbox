package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

type params struct {
	NetAddr  net.IP
	NetMask  net.IPMask
	Password string
	Ports    []uint
	Username string
}

func (p *params) Dump() string {
	return fmt.Sprintf("")
}

func NewParams() (*params, error) {
	var netAddr cmdtoolbox.IPParam
	var netMask cmdtoolbox.IPMaskParam
	var ports portsParam

	flag.Var(&netAddr, "addr", "IP address of the network")
	flag.Var(&netMask, "mask", "IP address of the network mask")
	password := flag.String("password", "", "password for the DVR")
	flag.Var(&ports, "port", "port number")
	username := flag.String("username", "admin", "username for the DVR")

	flag.Parse()

	return &params{
		NetAddr:  net.IP(netAddr),
		NetMask:  net.IPMask(netMask),
		Password: *password,
		Ports:    ports,
		Username: *username,
	}, nil
}

type portsParam []uint

func (param *portsParam) String() string {
	return "port parameters"
}

func (param *portsParam) Set(value string) error {
	v, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return err
	}

	*param = append(*param, uint(v))

	return nil
}
