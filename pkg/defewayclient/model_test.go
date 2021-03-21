package defewayclient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefewayJuan(t *testing.T) {
	t.Run("marshal to XML with empty RecSearch", func(t *testing.T) {
		expectedMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"></juan>`
		juan := DefewayJuan{}
		marshaled, err := juan.Marshal()

		require.NoError(t, err)
		require.Equal(t, expectedMarshaled, marshaled)
	})

	t.Run("marshal to XML with RecSearch", func(t *testing.T) {
		expectedMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"><recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0"></recsearch></juan>`
		juan := DefewayJuan{
			RecSearch: &DefewayRecSearch{
				Username:  "admin",
				Password:  "passwd",
				Channels:  uint16(3),
				Types:     uint16(15),
				Date:      "2019-01-01",
				BeginTime: "00:00:00",
				EndTime:   "23:59:59",
			},
		}
		marshaled, err := juan.Marshal()

		require.NoError(t, err)
		require.Equal(t, expectedMarshaled, marshaled)
	})

	t.Run("marshal to XML with DeviceInfo", func(t *testing.T) {
		expectedMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"><devinfo name="NVR" model="CS-580" serialnumber="AA000000000000" hwver="2.1.0" swver="2.5.2.10_22322230" reldatetime="2016/08/26 16:45" ip="192.168.1.1" httpport="80" clientport="8080" rip="123.123.123.123" rhttpport="60001" rclinetport="60001" camcnt="32"></devinfo></juan>`
		juan := DefewayJuan{
			DeviceInfo: &DefewayDeviceInfo{
				Name:             "NVR",
				Model:            "CS-580",
				SerialNumber:     "AA000000000000",
				HWVer:            "2.1.0",
				SWVer:            "2.5.2.10_22322230",
				RelDateTime:      "2016/08/26 16:45",
				IP:               "192.168.1.1",
				ClientPort:       8080,
				HTTPPort:         80,
				RemoteIP:         "123.123.123.123",
				RemoteHTTPPort:   60001,
				RemoteClientPort: 60001,
				CamCount:         32,
			},
		}

		marshaled, err := juan.Marshal()

		require.NoError(t, err)
		require.Equal(t, expectedMarshaled, marshaled)
	})

	t.Run("marshal to XML with EnvLoad", func(t *testing.T) {
		expectedMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"><envload usr="admin" pwd="p@ssw0rd" type="0" errno="0"></envload></juan>`
		juan := DefewayJuan{
			EnvLoad: &DefewayEnvLoad{
				Username: "admin",
				Password: "p@ssw0rd",
				Type:     0,
				ErrorNo:  0,
			},
		}

		marshaled, err := juan.Marshal()

		require.NoError(t, err)
		require.Equal(t, expectedMarshaled, marshaled)
	})
}

func TestUnmarshalJuan(t *testing.T) {
	t.Run("unmarshal XML from bytes with empty body", func(t *testing.T) {
		juanMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"></juan>`
		juan, err := UnmarshalJuan([]byte(juanMarshaled))

		require.NoError(t, err)
		require.Nil(t, juan.RecSearch)
		require.Nil(t, juan.EnvLoad)
		require.Nil(t, juan.DeviceInfo)

		validateJuan(t, juan)
	})

	t.Run("unmarshal XML from bytes with RecSearch with search results", func(t *testing.T) {
		juanMarshaled := `
		<juan ver="" squ="" dir="0" enc="0" errno="0">
			<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0">
				<s>0|1|3|8|1572887777|1572887780</s>
				<s>0|2|3|8|1572888888|1572888890</s>
			</recsearch>
		</juan>`
		juan, err := UnmarshalJuan([]byte(juanMarshaled))

		require.NoError(t, err)
		require.NotNil(t, juan.RecSearch)
		validateRecSearch(t, juan.RecSearch)
		validateJuan(t, juan)
	})

	t.Run("unmarshall XML from bytes with DevInfo and EnvLoad body", func(t *testing.T) {
		juanMarshaled := `
		<juan ver="" squ="" dir="0" enc="0" errno="0">
			<envload usr="admin" pwd="p@ssw0rd" type="0" errno="0"></envload>
			<devinfo name="NVR" model="CS-580" serialnumber="AA000000000000" hwver="2.1.0" swver="2.5.2.10_22322230" reldatetime="2016/08/26 16:45" ip="192.168.1.1" httpport="80" clientport="8080" rip="123.123.123.123" rhttpport="60001" rclinetport="60001" camcnt="32"></devinfo>
		</juan>`
		juan, err := UnmarshalJuan([]byte(juanMarshaled))

		require.NoError(t, err)
		require.NotNil(t, juan.EnvLoad)
		require.NotNil(t, juan.DeviceInfo)
		validateEnvLoad(t, juan.EnvLoad)
		validateDeviceInfo(t, juan.DeviceInfo)
	})

	t.Run("unmarshall XML from bytes with HDD with results", func(t *testing.T) {
		juanMarshaled := `
		<juan ver="" squ="" dir="0" errno="0">
		    <hdd errno="0" usr="admin" pwd="p@ssw0rd" action="0">
			    <d>Seagate 12345|5|2000|1000</d>
			</hdd>
		</juan>`
		juan, err := UnmarshalJuan([]byte(juanMarshaled))

		require.NoError(t, err)
		require.NotNil(t, juan.HDD)
		validateHDD(t, juan.HDD)
		validateJuan(t, juan)
	})
}

func TestRecordingMeta_GetFileName(t *testing.T) {
	t.Run("should return file name containing RecordingID, ChannelID and TypeID", func(t *testing.T) {
		rec := RecordingMeta{
			RecordingID: 1,
			ChannelID:   2,
			TypeID:      3,
		}

		fileName := rec.GetFileName()

		require.Equal(t, "1-2-3.flv", fileName)
	})
}

func TestRecordingMeta_GetFileShortName(t *testing.T) {
	t.Run("should return short file name containing RecordingID", func(t *testing.T) {
		rec := RecordingMeta{
			RecordingID: 1,
		}

		fileName := rec.GetFileShortName()

		require.Equal(t, "1.flv", fileName)
	})
}

func validateJuan(t *testing.T, juan *DefewayJuan) {
	require.Equal(t, "", juan.Version)
	require.Equal(t, "", juan.SQU)
	require.Equal(t, uint(0), juan.Direction)
	require.Equal(t, uint(0), juan.Enc)
	require.Equal(t, uint(0), juan.ErrorNo)
}

func validateRecSearch(t *testing.T, recSearch *DefewayRecSearch) {
	require.Equal(t, "admin", recSearch.Username)
	require.Equal(t, "passwd", recSearch.Password)
	require.Equal(t, uint16(3), recSearch.Channels)
	require.Equal(t, uint16(15), recSearch.Types)
	require.Equal(t, "2019-01-01", recSearch.Date)
	require.Equal(t, "00:00:00", recSearch.BeginTime)
	require.Equal(t, "23:59:59", recSearch.EndTime)
	require.Equal(t, uint(0), recSearch.SessionIdx)
	require.Equal(t, uint(0), recSearch.SessionCount)
	require.Equal(t, uint(0), recSearch.SessionTotal)

	validateSearchResults(t, recSearch.SearchResults)
}

func validateSearchResults(t *testing.T, searchResults []RecordingMeta) {
	require.Equal(t, 2, len(searchResults))

	res1 := searchResults[0]
	require.Equal(t, uint(1), res1.RecordingID)
	require.Equal(t, uint16(3), res1.ChannelID)
	require.Equal(t, uint16(8), res1.TypeID)
	require.Equal(t, uint64(1572887777), res1.StartTimestamp)
	require.Equal(t, uint64(1572887780), res1.EndTimestamp)

	res2 := searchResults[1]
	require.Equal(t, uint(2), res2.RecordingID)
	require.Equal(t, uint16(3), res2.ChannelID)
	require.Equal(t, uint16(8), res2.TypeID)
	require.Equal(t, uint64(1572888888), res2.StartTimestamp)
	require.Equal(t, uint64(1572888890), res2.EndTimestamp)
}

func validateEnvLoad(t *testing.T, del *DefewayEnvLoad) {
	require.Equal(t, "admin", del.Username)
	require.Equal(t, "p@ssw0rd", del.Password)
	require.Equal(t, uint8(0), del.Type)
	require.Equal(t, uint8(0), del.ErrorNo)
}

func validateDeviceInfo(t *testing.T, ddi *DefewayDeviceInfo) {
	require.Equal(t, "NVR", ddi.Name)
	require.Equal(t, "CS-580", ddi.Model)
	require.Equal(t, "AA000000000000", ddi.SerialNumber)
	require.Equal(t, "2.1.0", ddi.HWVer)
	require.Equal(t, "2.5.2.10_22322230", ddi.SWVer)
	require.Equal(t, "2016/08/26 16:45", ddi.RelDateTime)
	require.Equal(t, "192.168.1.1", ddi.IP)
	require.Equal(t, uint16(80), ddi.HTTPPort)
	require.Equal(t, uint16(8080), ddi.ClientPort)
	require.Equal(t, "123.123.123.123", ddi.RemoteIP)
	require.Equal(t, uint16(60001), ddi.RemoteHTTPPort)
	require.Equal(t, uint16(60001), ddi.RemoteClientPort)
	require.Equal(t, uint8(32), ddi.CamCount)
}

func validateHDD(t *testing.T, hdd *DefewayHDD) {
	require.Equal(t, "admin", hdd.Username)
	require.Equal(t, "p@ssw0rd", hdd.Password)
	require.Equal(t, uint8(0), hdd.Action)

	require.Equal(t, 1, len(hdd.Disks))

	disk1 := hdd.Disks[0]
	require.Equal(t, "Seagate 12345", disk1.Model)
	require.Equal(t, uint8(5), disk1.Status)
	require.Equal(t, uint64(2000), disk1.Capacity)
	require.Equal(t, uint64(1000), disk1.Used)
}
