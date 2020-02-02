package scanner

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sync"
)

type command struct {
	params ScannerParams
}

func NewCommand(params ScannerParams) *command {
	return &command{
		params: params,
	}
}

func (c *command) Run() error {

	var wg sync.WaitGroup

	netIP := c.params.NetAddr.Mask(c.params.NetMask)
	fmt.Println(netIP.String())

	netOnes, netBase := c.params.NetMask.Size()
	netSize := uint32(math.Pow(2, float64((netBase - netOnes))))
	ipStart := binary.BigEndian.Uint32(netIP)
	ipEnd := ipStart + netSize

	fmt.Printf("%+v  %+v  %+v  %+v\n", ipStart, netOnes, netBase, netSize)

	addrChan := make(chan string)

	wg.Add(1)
	go func(addrChan chan<- string) {
		defer wg.Done()
		defer close(addrChan)

		for ipCurr := ipStart; ipCurr < ipEnd; ipCurr++ {
			for _, port := range c.params.Ports {
				ip := make(net.IP, 4)
				binary.BigEndian.PutUint32(ip, ipCurr)
				addr := fmt.Sprintf("%s:%d", ip.String(), port)

				addrChan <- addr
			}
		}
	}(addrChan)

	for i := 0; i < c.params.Concurrent; i++ {
		wg.Add(1)
		go func(addrChan <-chan string) {
			defer wg.Done()

			for addr := range addrChan {
				fmt.Println(addr)
			}
		}(addrChan)
	}

	// for i := 0; i < c.params.Concurrent; i++ {
	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		if err := c.scan(); err != nil {
	// 			log.Prinln(err)
	// 		}
	// 	}()
	// }

	wg.Wait()
	return nil
}
