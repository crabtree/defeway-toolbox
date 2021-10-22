package downloader

import (
	"io"
	"log"
	"sync"

	dc "github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

type RecordingsClient interface {
	Fetch(fetchParams dc.RecordingsFetchParams) ([]dc.RecordingMeta, error)
	Download(recMeta dc.RecordingMeta, dst io.Writer, isPreview bool) error
}

type command struct {
	client RecordingsClient
	params DownloaderParams
}

func NewCommand(client RecordingsClient, params DownloaderParams) *command {
	return &command{
		client: client,
		params: params,
	}
}

func (c *command) Run() error {
	var wg sync.WaitGroup

	recsChan, err := c.fetch()
	if err != nil {
		return err
	}

	for i := 0; i < c.params.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err = c.process(recsChan); err != nil {
				log.Println(err)
			}
		}()
	}

	wg.Wait()
	return nil
}
