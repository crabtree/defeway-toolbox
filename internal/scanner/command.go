package scanner

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
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
			c.getClientConfig(addr))

		info, err := client.Fetch()
		if err != nil {
			log.Println(err)
			continue
		}

		payload := fmt.Sprintf(`<a href="http://%s">http://%s</a>`, addr, addr)
		fileNameBase := fmt.Sprintf("%s.html", strings.ReplaceAll(addr, ":", "-"))

		if info.EnvLoad.ErrorNo != 0 {
			logFilePath := path.Join(c.params.LogDir, fmt.Sprintf("s-%s", fileNameBase))
			writeLog(logFilePath, payload)
			log.Printf("Found device http://%s, with env error %d\n", addr, info.EnvLoad.ErrorNo)
			continue
		}

		log.Printf("Found device http://%s\n", addr)

		logFilePath := path.Join(c.params.LogDir, fileNameBase)
		infoSerialized, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			writeLog(logFilePath, payload)
			continue
		} else {
			payload += fmt.Sprintf(`<br><pre>%s</pre>`, string(infoSerialized))
			writeLog(logFilePath, payload)
		}

		dstPath := path.Join(c.params.LogDir, strings.ReplaceAll(addr, ":", "-"))
		if err := cmdtoolbox.EnsureDir(dstPath); err != nil {
			return err
		}
		var ch uint8 = 0
		snapshotClient := defewayclient.NewSnapshotClient(c.getClientConfig(addr))
		for ; ch < info.DeviceInfo.CamCount; ch++ {
			fp := path.Join(dstPath, fmt.Sprintf("ch-%d.jpg", ch))
			err := c.fetchSnapshotForCh(snapshotClient, int(ch), fp)
			if err != nil {
				log.Printf("Error: %s\n", err)
				os.Remove(fp)
			}
		}
	}

	return nil
}

func (c *command) getClientConfig(addr string) defewayclient.DefewayClientConfig {
	return defewayclient.DefewayClientConfig{
		Address:  addr,
		Username: c.params.Username,
		Password: c.params.Password,
		HTTPClientConfig: defewayclient.HTTPClientConfig{
			Timeout:           c.params.Timeout,
			TLSSkipVerify:     c.params.TLSSkipVerify,
			DisableKeepAlives: true,
		},
	}
}

func (c *command) fetchSnapshotForCh(client *defewayclient.SnapshotClient, ch int, dstPath string) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		err := dst.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return client.Fetch(int(ch), dst)
}

func writeLog(logFilePath, payload string) {
	if err := ioutil.WriteFile(logFilePath, []byte(payload), 0644); err != nil {
		log.Println(err)
	}
}
