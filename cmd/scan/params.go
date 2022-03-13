package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

type params struct {
	Concurrent    int
	LogDir        string
	NetAddr       net.IP
	NetMask       net.IPMask
	Password      string
	Ports         []uint
	Timeout       time.Duration
	TLSSkipVerify bool
	Username      string
	WithSnapshots bool
}

func (p *params) Dump() string {
	return fmt.Sprintf("Concurrent=%d LogDir=%s NetAddr=%s NetMask=%s Password=%s Ports=%d Timeout=%d TLSSkipVerify=%t Username=%s WithSnapshots=%t",
		p.Concurrent, p.LogDir, p.NetAddr, p.NetMask, p.Password, p.Ports, p.Timeout, p.TLSSkipVerify, p.Username, p.WithSnapshots)
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
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "disables the TLS certificate verification")
	timeout := flag.Duration("timeout", 5*time.Second, "sets the client timeout")
	username := flag.String("username", "admin", "username for the DVR")
	withSnapshots := flag.Bool("with-snapshots", false, "fetch channels snapshots of discovered device")

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
		Concurrent:    *concurrent,
		LogDir:        *logDir,
		NetAddr:       net.IP(netAddr),
		NetMask:       net.IPMask(netMask),
		Password:      *password,
		Ports:         ports,
		TLSSkipVerify: *tlsSkipVerify,
		Timeout:       *timeout,
		Username:      *username,
		WithSnapshots: *withSnapshots,
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
