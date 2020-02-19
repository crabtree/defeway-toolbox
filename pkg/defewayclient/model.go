package defewayclient

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type DefewayJuan struct {
	XMLName    xml.Name           `xml:"juan"`
	Version    string             `xml:"ver,attr"`
	SQU        string             `xml:"squ,attr"`
	Direction  uint               `xml:"dir,attr"`
	Enc        uint               `xml:"enc,attr"`
	ErrorNo    uint               `xml:"errno,attr"`
	RecSearch  *DefewayRecSearch  `xml:"recsearch,omitempty"`
	DeviceInfo *DefewayDeviceInfo `xml:"devinfo,omitempty"`
	EnvLoad    *DefewayEnvLoad    `xml:"envload,omitempty"`
}

func (dj *DefewayJuan) Marshal() (string, error) {
	b, err := xml.Marshal(dj)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func NewForRecSearch(recSearch DefewayRecSearch) *DefewayJuan {
	return &DefewayJuan{
		RecSearch: &recSearch,
	}
}

func UnmarshalJuanForRecSearch(data []byte) (*DefewayJuan, error) {
	dj := &DefewayJuan{}
	err := xml.Unmarshal(data, dj)
	if err != nil {
		return nil, err
	}

	return dj, nil
}

type DefewayRecSearch struct {
	Username      string          `xml:"usr,attr"`
	Password      string          `xml:"pwd,attr"`
	Channels      uint16          `xml:"channels,attr"`
	Types         uint16          `xml:"types,attr"`
	Date          string          `xml:"date,attr"`
	BeginTime     string          `xml:"begin,attr"`
	EndTime       string          `xml:"end,attr"`
	SessionIdx    uint            `xml:"session_index,attr"`
	SessionCount  uint            `xml:"session_count,attr"`
	SessionTotal  uint            `xml:"session_total,attr"`
	SearchResults []RecordingMeta `xml:"s,omitempty"`
}

type RecordingMeta struct {
	RecordingID    uint
	ChannelID      uint16
	TypeID         uint16
	StartTimestamp uint64
	EndTimestamp   uint64
}

func (s *RecordingMeta) GetFileShortName() string {
	return fmt.Sprintf("%d.flv", s.RecordingID)
}

func (s *RecordingMeta) GetFileName() string {
	return fmt.Sprintf("%d-%d-%d.flv", s.RecordingID, s.ChannelID, s.TypeID)
}

func (s *RecordingMeta) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var val string
	if err := d.DecodeElement(&val, &start); err != nil {
		return err
	}

	valSplitted := strings.Split(val, "|")

	recID, err := strconv.ParseUint(valSplitted[1], 10, 32)
	if err != nil {
		return err
	}

	channelID, err := strconv.ParseUint(valSplitted[2], 10, 16)
	if err != nil {
		return err
	}

	typeID, err := strconv.ParseUint(valSplitted[3], 10, 16)
	if err != nil {
		return err
	}

	startTimestamp, err := strconv.ParseUint(valSplitted[4], 10, 64)
	if err != nil {
		return err
	}

	endTimestamp, err := strconv.ParseUint(valSplitted[5], 10, 64)
	if err != nil {
		return err
	}

	*s = RecordingMeta{
		RecordingID:    uint(recID),
		ChannelID:      uint16(channelID),
		TypeID:         uint16(typeID),
		StartTimestamp: startTimestamp,
		EndTimestamp:   endTimestamp,
	}

	return nil
}

type DefewayDeviceInfo struct {
	Name             string `xml:"name,attr"`
	Model            string `xml:"model,attr"`
	SerialNumber     string `xml:"serialnumber,attr"`
	HWVer            string `xml:"hwver,attr"`
	SWVer            string `xml:"swver,attr"`
	RelDateTime      string `xml:"reldatetime,attr"`
	IP               string `xml:"ip,attr"`
	HTTPPort         uint16 `xml:"httpport,attr"`
	ClientPort       uint16 `xml:"clientport,attr"`
	RemoteIP         string `xml:"rip,attr"`
	RemoteHTTPPort   uint16 `xml:"rhttpport,attr"`
	RemoteClientPort uint16 `xml:"rclinetport,attr"`
	CamCount         uint8  `xml:"camcnt,attr"`
}

func NewForDeviceInfo(envLoad DefewayEnvLoad, devInfo DefewayDeviceInfo) *DefewayJuan {
	return &DefewayJuan{
		DeviceInfo: &devInfo,
		EnvLoad:    &envLoad,
	}
}

func UnmarshalJuanForDeviceInfo(data []byte) (*DefewayJuan, error) {
	dj := &DefewayJuan{}
	err := xml.Unmarshal(data, dj)
	if err != nil {
		return nil, err
	}

	return dj, nil
}

type DefewayEnvLoad struct {
	Username string `xml:"usr,attr"`
	Password string `xml:"pwd,attr"`
	Type     uint8  `xml:"type,attr"`
	ErrorNo  uint8  `xml:"errno,attr"`
}
