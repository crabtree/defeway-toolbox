package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

func main() {
	params, err := NewParams()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(params.Dump())

	client := defewayclient.NewDefewayClient(
		fmt.Sprintf("%s:%d", params.Address, params.Port),
		params.Username,
		params.Password)
	recMgr := defewayclient.NewRecordingsManager(client)
	recordings, err := recMgr.Fetch(paramsToRecordingsFetchParams(params))
	if err != nil {
		log.Fatal(err)
	}

	outDir := path.Join(params.OutputDir, params.Address.String())
	if err = ensureRecordingsDir(outDir); err != nil {
		log.Fatal(err)
	}

	recordingsChan := make(chan defewayclient.RecordingMeta, len(recordings))

	var wg sync.WaitGroup

	wg.Add(1)
	go addRecordingsToChannel(&wg, recordings, recordingsChan)

	for i := 0; i < params.Concurrent; i++ {
		wg.Add(1)
		go downloadRecording(&wg, recordingsChan, outDir, params.Overwrite, recMgr)
	}

	log.Printf("Waiting to process the recordings (%d)\n", len(recordings))
	wg.Wait()
}

func paramsToRecordingsFetchParams(params *params) defewayclient.RecordingsFetchParams {
	return defewayclient.RecordingsFetchParams{
		Channels:       params.Recordings.Channels,
		Date:           params.Recordings.Date,
		EndTime:        params.Recordings.EndTime,
		RecordingTypes: params.Recordings.RecordingTypes,
		StartTime:      params.Recordings.StartTime,
	}
}

func addRecordingsToChannel(wg *sync.WaitGroup, recordings []defewayclient.RecordingMeta, recordingsChan chan<- defewayclient.RecordingMeta) {
	defer wg.Done()

	for _, rec := range recordings {
		recordingsChan <- rec
	}

	close(recordingsChan)
}

func ensureRecordingsDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}

	return err
}

func downloadRecording(
	wg *sync.WaitGroup,
	recordingsChan <-chan defewayclient.RecordingMeta,
	outDir string,
	overwrite bool,
	recMgr *defewayclient.RecordingsManager) {

	defer wg.Done()

	for rec := range recordingsChan {
		dstPath := path.Join(outDir, fmt.Sprintf("%d.flv", rec.RecordingID))
		exists, err := fileExists(dstPath)
		if err != nil {
			log.Println(err)
			return
		}

		if exists && !overwrite {
			log.Printf("File %s already exists\n", dstPath)
			continue
		}

		log.Printf("Downloading %d into %s\n", rec.RecordingID, dstPath)

		dst, err := os.Create(dstPath)
		if err != nil {
			log.Println(err)
			return
		}
		defer dst.Close()

		if err = recMgr.Download(rec, dst); err != nil {
			log.Println(err)
			return
		}
	}
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
