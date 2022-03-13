package scanner

import (
	"net"
	"time"
)

type ScannerParams struct {
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
