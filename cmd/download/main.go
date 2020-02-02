package main

import (
	"fmt"
	"log"
	"path"

	"github.com/crabtree/defeway-toolbox/internal/downloader"
	"github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

func main() {
	params, err := NewParams()
	dieOnError(err)

	log.Println(params.Dump())

	client := defewayclient.NewRecordingsClient(
		fmt.Sprintf("%s:%d", params.Client.Address, params.Client.Port),
		params.Client.Username,
		params.Client.Password)

	command := downloader.NewCommand(
		client,
		paramsToCommandParams(params))

	err = command.Run()
	dieOnError(err)
}

func paramsToCommandParams(params *params) downloader.DownloaderParams {
	return downloader.DownloaderParams{
		Channels:   params.Recordings.Channels,
		Concurrent: params.Downloads.Concurrent,
		Date:       params.Recordings.Date,
		EndTime:    params.Recordings.EndTime,
		Overwrite:  params.Downloads.Overwrite,
		OutputDir: path.Join(
			params.Downloads.OutputDir,
			params.Recordings.Date.Format("2006-01-02"),
			params.Client.Address.String()),
		RecordingTypes: params.Recordings.RecordingTypes,
		StartTime:      params.Recordings.StartTime,
	}
}

func dieOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
