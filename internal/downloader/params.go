package downloader

import (
	"time"

	dc "github.com/crabtree/defeway-toolbox/pkg/defewayclient"
)

type DownloaderParams struct {
	Channels       uint16
	Concurrent     int
	Date           time.Time
	EndTime        time.Time
	Overwrite      bool
	OutputDir      string
	RecordingTypes uint16
	StartTime      time.Time
}

func (dp *DownloaderParams) ToRecordingsFetchParams() dc.RecordingsFetchParams {
	return dc.RecordingsFetchParams{
		Channels:       dp.Channels,
		Date:           dp.Date,
		EndTime:        dp.EndTime,
		RecordingTypes: dp.RecordingTypes,
		StartTime:      dp.StartTime,
	}
}
