package scanner

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"path"
	"strings"
	"sync"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
	"github.com/crabtree/defeway-toolbox/pkg/defewayclient"
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

	if err := cmdtoolbox.EnsureDir(c.params.LogDir); err != nil {
		return err
	}

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

	netIP := c.params.NetAddr.Mask(c.params.NetMask)
	netOnes, netBase := c.params.NetMask.Size()
	netSize := uint32(math.Pow(2, float64((netBase - netOnes))))
	ipStart := binary.BigEndian.Uint32(netIP)
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

		payload := fmt.Sprintf(`<a href="http://%s">http://%s</a>`, addr, addr)
		fileNameBase := fmt.Sprintf("%s.html", strings.ReplaceAll(addr, ":", "-"))

		if info.EnvLoad.ErrorNo != 0 {
			logFilePath := path.Join(c.params.LogDir, fmt.Sprintf("s-%s", fileNameBase))
			if err = ioutil.WriteFile(logFilePath, []byte(payload), 0644); err != nil {
				log.Println(err)
			}

			log.Printf("Found device http://%s, with env error %d\n", addr, info.EnvLoad.ErrorNo)
			continue
		}

		logFilePath := path.Join(c.params.LogDir, fileNameBase)
		if err = ioutil.WriteFile(logFilePath, []byte(payload), 0644); err != nil {
			log.Println(err)
		}

		log.Printf("Found device http://%s\n", addr)
	}

	return nil
}
