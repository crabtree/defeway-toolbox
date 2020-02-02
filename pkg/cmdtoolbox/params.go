package cmdtoolbox

import (
	"fmt"
	"net"
)

type IPParam net.IP

func (param *IPParam) String() string {
	return "IP address parameter"
}

func (param *IPParam) Set(value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return fmt.Errorf("the value %s is not a valid IP address", value)
	}

	*param = IPParam(ip)
	return nil
}

type IPMaskParam net.IPMask

func (param *IPMaskParam) String() string {
	return "IP mask parameter"
}

func (param *IPMaskParam) Set(value string) error {
	ip := net.ParseIP(value)
	if ip == nil {
		return fmt.Errorf("the value %s is not a valid IP address", value)
	}

	*param = IPMaskParam(net.IPMask(ip))
	return nil
}
