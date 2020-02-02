package main

import (
	"fmt"
	"log"
	"path"

	"github.com/crabtree/defeway-toolbox/internal/downloader"
	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
	"github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

func main() {
	params, err := NewParams()
	cmdtoolbox.DieOnError(err)

	log.Println(params.Dump())

	client := defewayclient.NewRecordingsClient(
		fmt.Sprintf("%s:%d", params.Client.Address, params.Client.Port),
		params.Client.Username,
		params.Client.Password)

	command := downloader.NewCommand(
		client,
		paramsToCommandParams(params))

	err = command.Run()
	cmdtoolbox.DieOnError(err)
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
			params.Client.Address.String(),
			params.Recordings.Date.Format("2006-01-02")),
		RecordingTypes: params.Recordings.RecordingTypes,
		StartTime:      params.Recordings.StartTime,
	}
}
