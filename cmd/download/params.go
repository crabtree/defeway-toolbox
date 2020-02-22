package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/crabtree/defeway-toolbox/pkg/cmdtoolbox"
)

type clientParams struct {
	Address           net.IP
	DisableKeepAlives bool
	Password          string
	Port              uint
	Timeout           time.Duration
	TLSSkipVerify     bool
	Username          string
}

func (p *clientParams) Dump() string {
	return fmt.Sprintf("Address=%s DisableKeepAlives=%t Password=%s Port=%d Timeout=%d TLSSkipVerify=%t Username=%s",
		p.Address, p.DisableKeepAlives, p.Password, p.Port, p.Timeout, p.TLSSkipVerify, p.Username)
}

type downloadsParams struct {
	Concurrent int
	OutputDir  string
	Overwrite  bool
}

func (p *downloadsParams) Dump() string {
	return fmt.Sprintf("Concurrent=%d Output=%s Overwrite=%t",
		p.Concurrent, p.OutputDir, p.Overwrite)
}

type recordingsParams struct {
	Channels       uint16
	Date           time.Time
	EndTime        time.Time
	RecordingTypes uint16
	StartTime      time.Time
}

func (p *recordingsParams) Dump() string {
	return fmt.Sprintf("Channels=%d Date=%s EndTime=%s RecordingTypes=%d StartTime=%s",
		p.Channels, p.Date.Format("2006-01-02"), p.EndTime.Format("15:04:05"), p.RecordingTypes, p.StartTime.Format("15:04:05"))
}

type params struct {
	Client     *clientParams
	Downloads  *downloadsParams
	Recordings *recordingsParams
}

func (p *params) Dump() string {
	return fmt.Sprintf("%s %s %s",
		p.Client.Dump(), p.Downloads.Dump(), p.Recordings.Dump())
}

func NewParams() (*params, error) {
	var address cmdtoolbox.IPParam
	var channels channelsParam
	var date dateParam
	var endTime timeParam
	var startTime timeParam
	var types recordingTypesParam

	flag.Var(&address, "addr", "IP address of the DVR")
	flag.Var(&channels, "chan", "channel id")
	concurrent := flag.Int("concurrent", 1, "sets the number of concurrent workers")
	flag.Var(&date, "date", "specify date in format YYYY-MM-DD (eg. 2019-01-01)")
	disableKeepAlives := flag.Bool("no-keep-alives", false, "disables the keep alives connections")
	flag.Var(&endTime, "end", "recording end time")
	outputDir := flag.String("output", "", "path to the downloads directory")
	overwrite := flag.Bool("overwrite", false, "overwrite existing files")
	password := flag.String("password", "", "password for the DVR")
	port := flag.Int("port", 60001, "sets the port to the DVR")
	flag.Var(&startTime, "start", "recording start time")
	tlsSkipVerify := flag.Bool("tls-skip-verify", false, "disables the TLS certificate verification")
	timeout := flag.Duration("timeout", 5*time.Second, "sets the client timeout")
	flag.Var(&types, "type", "recording type")
	username := flag.String("username", "admin", "username for the DVR")

	flag.Parse()

	if address == nil {
		return nil, fmt.Errorf("specify IP address")
	}

	if channels == 0 {
		return nil, fmt.Errorf("specify at least one channel id")
	}

	if time.Time(endTime).IsZero() {
		endTime = timeParam(time.Date(0, 0, 0, 23, 59, 59, 999999999, time.UTC))
	}

	if time.Time(date).IsZero() {
		date = dateParam(time.Now().Add(-24 * time.Hour))
	}

	if outputDir == nil || *outputDir == "" {
		return nil, fmt.Errorf("specify downloads directory")
	}

	if types == 0 {
		return nil, fmt.Errorf("specify at least one recording type")
	}

	return &params{
		Client: &clientParams{
			Address:           net.IP(address),
			DisableKeepAlives: *disableKeepAlives,
			Password:          *password,
			Port:              uint(*port),
			TLSSkipVerify:     *tlsSkipVerify,
			Timeout:           *timeout,
			Username:          *username,
		},
		Downloads: &downloadsParams{
			Concurrent: *concurrent,
			OutputDir:  *outputDir,
			Overwrite:  *overwrite,
		},
		Recordings: &recordingsParams{
			Channels:       uint16(channels),
			Date:           time.Time(date),
			EndTime:        time.Time(endTime),
			RecordingTypes: uint16(types),
			StartTime:      time.Time(startTime),
		},
	}, nil
}

type channelsParam uint16

func (c *channelsParam) String() string {
	return "cameras parameters"
}

func (c *channelsParam) Set(value string) error {
	v, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return err
	}

	*c = *c | channelsParam(math.Pow(float64(2), float64(v-1)))

	return nil
}

type recordingTypesParam uint16

func (rt *recordingTypesParam) String() string {
	return "recording types"
}

func (rt *recordingTypesParam) Set(value string) error {
	v, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		return err
	}

	*rt = *rt | recordingTypesParam(math.Pow(float64(2), float64(v-1)))

	return nil
}

type dateParam time.Time

func (dp *dateParam) String() string {
	return "date parameter"
}

func (dp *dateParam) Set(value string) error {
	v, err := time.Parse("2006-01-02", value)
	if err != nil {
		return err
	}

	*dp = dateParam(v)

	return nil
}

type timeParam time.Time

func (tp *timeParam) String() string {
	return "time parameter"
}

func (tp *timeParam) Set(value string) error {
	v, err := time.Parse("15:04:05", value)
	if err != nil {
		return err
	}

	*tp = timeParam(v)

	return nil
}
