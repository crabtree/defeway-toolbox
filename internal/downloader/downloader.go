package downloader

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"

	dc "github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

type RecordingsManager interface {
	Fetch(fetchParams dc.RecordingsFetchParams) ([]dc.RecordingMeta, error)
	Download(recMeta dc.RecordingMeta, dst io.Writer) error
}

type cmd struct {
	recMgr  RecordingsManager
	recChan chan dc.RecordingMeta
}

func NewCmd(recMgr RecordingsManager, recsChan chan dc.RecordingMeta) *cmd {
	return &cmd{
		recMgr:  recMgr,
		recChan: recsChan,
	}
}

func (c *cmd) FetchRecordings(params dc.RecordingsFetchParams) error {
	defer close(c.recChan)

	recordings, err := c.recMgr.Fetch(params)
	if err != nil {
		return err
	}

	for _, rec := range recordings {
		c.recChan <- rec
	}

	return nil
}

func (c *cmd) ProcessRecordings(outDir string, overwrite bool) error {
	for rec := range c.recChan {
		dstPath := path.Join(outDir, getFileName(rec))
		exists, err := fileExists(dstPath)
		if err != nil {
			log.Println(err)
			continue
		}

		if exists && !overwrite {
			log.Printf("File %s already exists\n", dstPath)
			continue
		}

		log.Printf("Downloading %d into %s\n", rec.RecordingID, dstPath)

		err = c.downloadFile(dstPath, rec)
		if err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}

func (c *cmd) downloadFile(dstPath string, recMeta dc.RecordingMeta) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if err = c.recMgr.Download(recMeta, dst); err != nil {
		return err
	}

	return nil
}

func getFileName(rec dc.RecordingMeta) string {
	return fmt.Sprintf("%d.%d.%d.flv", rec.RecordingID, rec.ChannelID, rec.TypeID)
}

func fileExists(dstPath string) (bool, error) {
	_, err := os.Stat(dstPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
