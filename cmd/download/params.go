package main

import (
	"flag"
	"fmt"
	"math"
	"net"
	"strconv"
	"time"
)

type params struct {
	Address    net.IP
	Concurrent int
	OutputDir  string
	Overwrite  bool
	Password   string
	Port       uint
	Recordings *recordingsParams
	Username   string
}

func (p *params) Dump() string {
	return fmt.Sprintf("Address=%s Concurrent=%d Output=%s Overwrite=%t Password=%s Port=%d Username=%s %s",
		p.Address, p.Concurrent, p.OutputDir, p.Overwrite, p.Password, p.Port, p.Username, p.Recordings.Dump())
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

func NewParams() (*params, error) {
	var address ipParam
	var channels channelsParam
	var date dateParam
	var endTime timeParam
	var startTime timeParam
	var types recordingTypesParam

	flag.Var(&address, "addr", "IP address of the DVR")
	flag.Var(&channels, "chan", "channel id")
	concurrent := flag.Int("concurrent", 1, "sets the number of concurrent workers")
	flag.Var(&date, "date", "specify date in format YYYY-MM-DD (eg. 2019-01-01)")
	flag.Var(&endTime, "end", "recording strat time")
	outputDir := flag.String("output", "", "path to the downloads directory")
	overwrite := flag.Bool("overwrite", false, "overwrite existing files")
	password := flag.String("password", "", "password for the DVR")
	port := flag.Int("port", 60001, "sets the port to the DVR")
	flag.Var(&startTime, "start", "recording strat time")
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
		Address:    net.IP(address),
		Concurrent: *concurrent,
		OutputDir:  *outputDir,
		Overwrite:  *overwrite,
		Password:   *password,
		Port:       uint(*port),
		Recordings: &recordingsParams{
			Channels:       uint16(channels),
			Date:           time.Time(date),
			EndTime:        time.Time(endTime),
			RecordingTypes: uint16(types),
			StartTime:      time.Time(startTime),
		},
		Username: *username,
	}, nil
}

type ipParam net.IP

func (ip *ipParam) String() string {
	return "IP address parameter"
}

func (ip *ipParam) Set(value string) error {
	_ip := net.ParseIP(value)
	if _ip == nil {
		return fmt.Errorf("the value %s is not a valid IP address", value)
	}

	*ip = ipParam(_ip)
	return nil
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
