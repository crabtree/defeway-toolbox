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

	ipArr := []byte(ip)

	*param = IPMaskParam(net.IPv4Mask(ipArr[12], ipArr[13], ipArr[14], ipArr[15]))
	return nil
}
