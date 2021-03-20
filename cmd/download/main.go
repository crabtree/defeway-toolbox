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
		paramsToClientConfig(params),
		paramsToDownloadClientConfig(params))

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
		InputFile:  params.Recordings.InputFile,
		Overwrite:  params.Downloads.Overwrite,
		OutputDir: path.Join(
			params.Downloads.OutputDir,
			fmt.Sprintf("%s-%d", params.Client.Address.String(), params.Client.Port),
			params.Recordings.Date.Format("2006-01-02")),
		RecordingTypes: params.Recordings.RecordingTypes,
		StartTime:      params.Recordings.StartTime,
	}
}

func paramsToClientConfig(params *params) defewayclient.DefewayClientConfig {
	return defewayclient.DefewayClientConfig{
		Address:  fmt.Sprintf("%s:%d", params.Client.Address, params.Client.Port),
		Username: params.Client.Username,
		Password: params.Client.Password,
		HTTPClientConfig: defewayclient.HTTPClientConfig{
			DisableKeepAlives: params.Client.DisableKeepAlives,
			TLSSkipVerify:     params.Client.TLSSkipVerify,
			Timeout:           params.Client.Timeout,
		},
	}
}

func paramsToDownloadClientConfig(params *params) defewayclient.DefewayClientConfig {
	cfg := paramsToClientConfig(params)
	cfg.Timeout = 0

	return cfg
}
