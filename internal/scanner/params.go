package scanner

import "net"

type ScannerParams struct {
	Concurrent int
	NetAddr    net.IP
	NetMask    net.IPMask
	Password   string
	Ports      []uint
	Username   string
}
