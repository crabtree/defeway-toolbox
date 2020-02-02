package downloader

import (
	"log"
	"os"
	"path"

	dc "github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

func (c *command) fetch() (<-chan dc.RecordingMeta, error) {
	recordings, err := c.client.Fetch(c.params.ToRecordingsFetchParams())
	if err != nil {
		return nil, err
	}

	recordingsChan := make(chan dc.RecordingMeta, len(recordings))
	defer close(recordingsChan)

	for _, rec := range recordings {
		recordingsChan <- rec
	}

	return recordingsChan, nil
}

func (c *command) process(recsChan <-chan dc.RecordingMeta) error {
	if err := ensureRecordingsDir(c.params.OutputDir); err != nil {
		return err
	}

	for recMeta := range recsChan {
		dstPath := path.Join(c.params.OutputDir, recMeta.GetFileName())
		exists, err := fileExists(dstPath)
		if err != nil {
			log.Println(err)
			continue
		}

		if exists && !c.params.Overwrite {
			log.Printf("File %s already exists\n", dstPath)
			continue
		}

		log.Printf("Downloading %d into %s\n", recMeta.RecordingID, dstPath)

		if err = c.download(dstPath, recMeta); err != nil {
			log.Println(err)
			continue
		}
	}

	return nil
}

func (c *command) download(dstPath string, recMeta dc.RecordingMeta) error {
	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func() {
		err := dst.Close()
		log.Println(err)
	}()

	if err = c.client.Download(recMeta, dst); err != nil {
		return err
	}

	return nil
}

func ensureRecordingsDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}

	return err
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
