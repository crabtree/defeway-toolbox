package defewayclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_DeviceInfoClient_Fetch(t *testing.T) {
	t.Run("returns error when error occures during http call", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte{})
		}))
		defer server.Close()

		c := &DeviceInfoClient{fixClient(server.Client(), "invalid-address")}

		_, err := c.Fetch()

		require.Contains(t, err.Error(), "no such host")
	})

	t.Run("returns error when max retry reached because of error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			juanMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="1"></juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		c := &DeviceInfoClient{fixClient(server.Client(), server.URL[7:])}

		_, err := c.Fetch()

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of invalid response format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(""))
		}))
		defer server.Close()

		c := &DeviceInfoClient{fixClient(server.Client(), server.URL[7:])}

		_, err := c.Fetch()

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of empty devinfo response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/gw.cgi`, req.URL.Path)
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<envload usr="admin" pwd="p@ssw0rd" type="0" errno="0"></envload>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		c := &DeviceInfoClient{fixClient(server.Client(), server.URL[7:])}

		_, err := c.Fetch()

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns device info", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/gw.cgi`, req.URL.Path)
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<envload usr="admin" pwd="p@ssw0rd" type="0" errno="0"></envload>
				<devinfo name="NVR" model="CS-580" serialnumber="AA000000000000" hwver="2.1.0" swver="2.5.2.10_22322230" reldatetime="2016/08/26 16:45" ip="192.168.1.1" httpport="80" clientport="8080" rip="123.123.123.123" rhttpport="60001" rclinetport="60001" camcnt="32"></devinfo>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		c := &DeviceInfoClient{fixClient(server.Client(), server.URL[7:])}

		juan, err := c.Fetch()

		require.NoError(t, err)
		require.NotNil(t, juan.EnvLoad)
		require.NotNil(t, juan.DeviceInfo)
	})
}
