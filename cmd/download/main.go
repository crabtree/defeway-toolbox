package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"

	"github.com/crabtree/defeway-toolbox/internal/downloader"
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

	recordingsChan := make(chan defewayclient.RecordingMeta)

	cmd := downloader.NewCmd(recMgr, recordingsChan)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()
		err := cmd.FetchRecordings(paramsToRecordingsFetchParams(params))
		if err != nil {
			log.Fatal(err)
		}
	}()

	outDir := path.Join(params.OutputDir, params.Address.String())
	if err = ensureRecordingsDir(outDir); err != nil {
		log.Fatal(err)
	}

	for i := 0; i < params.Concurrent; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := cmd.ProcessRecordings(outDir, params.Overwrite)
			if err != nil {
				log.Println(err)
			}
		}()
	}

	log.Println("Waiting to process the recordings")
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

func ensureRecordingsDir(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dirPath, 0755)
		return err
	}

	return err
}
