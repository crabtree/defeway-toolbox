package scanner

import (
	"encoding/binary"
	"fmt"
	"github.com/crabtree/defeway-toolbox/pkg/defewayclient"
	"log"
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
	addrChan := make(chan string, 100)

	wg.Add(1)
	go func(addrChan chan<- string) {
		defer wg.Done()
		c.prepareAddresses(addrChan)
	}(addrChan)

	for i := 0; i < c.params.Concurrent; i++ {
		wg.Add(1)
		go func(addrChan <-chan string) {
			defer wg.Done()
			if err := c.scan(addrChan); err != nil {
				log.Println(err)
			}
		}(addrChan)
	}

	wg.Wait()
	return nil
}

func (c *command) prepareAddresses(addrChan chan<- string) {
	defer close(addrChan)

	netOnes, netBase := c.params.NetMask.Size()
	netSize := uint32(math.Pow(2, float64((netBase - netOnes))))
	ipStart := binary.BigEndian.Uint32(c.params.NetAddr.To4())
	ipEnd := ipStart + netSize

	for ipCurr := ipStart; ipCurr < ipEnd; ipCurr++ {
		for _, port := range c.params.Ports {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, ipCurr)
			addr := fmt.Sprintf("%s:%d", ip.String(), port)

			addrChan <- addr
		}
	}
}

func (c *command) scan(addrChan <-chan string) error {
	for addr := range addrChan {
		client := defewayclient.NewDeviceInfoClient(
			addr,
			c.params.Username,
			c.params.Password)

		info, err := client.Fetch()
		if err != nil {
			log.Println(err)
			continue
		}

		if info.EnvLoad.ErrorNo != 0 {
			log.Printf("Found device http://%s, with env error %d\n", addr, info.EnvLoad.ErrorNo)
			continue
		}

		log.Printf("Found device http://%s\n", addr)
	}

	return nil
}
