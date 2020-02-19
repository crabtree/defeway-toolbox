package scanner

import (
	"net"
)

type ScannerParams struct {
	Concurrent int
	LogDir     string
	NetAddr    net.IP
	NetMask    net.IPMask
	Password   string
	Ports      []uint
	Username   string
}
