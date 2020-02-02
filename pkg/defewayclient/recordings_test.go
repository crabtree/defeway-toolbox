package defewayclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RecordingsClient_Fetch(t *testing.T) {
	t.Run("returns error when error occures during http call", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte{})
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			FetchClient: server.Client(),
			Address:     "invalid-address",
		}}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Contains(t, err.Error(), "no such host")
	})

	t.Run("returns error when max retry reached because of no recordings found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0">
				</recsearch>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			FetchClient: server.Client(),
			Address:     server.URL[7:],
		}}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of error response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="1">
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			FetchClient: server.Client(),
			Address:     server.URL[7:],
		}}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns error when max retry reached because of invalid response format", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(""))
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			FetchClient: server.Client(),
			Address:     server.URL[7:],
		}}
		fetchParams := RecordingsFetchParams{}

		_, err := rm.Fetch(fetchParams)

		require.Equal(t, "max retry count reached", err.Error())
	})

	t.Run("returns slice with recordings", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/gw.cgi`, req.URL.Path)
			juanMarshaled := `
			<juan ver="" squ="" dir="0" enc="0" errno="0">
				<recsearch usr="admin" pwd="passwd" channels="3" types="15" date="2019-01-01" begin="00:00:00" end="23:59:59" session_index="0" session_count="0" session_total="0">
					<s>0|1|3|8|1572887777|1572887780</s>
					<s>0|2|3|8|1572888888|1572888890</s>
				</recsearch>
			</juan>`
			rw.Write([]byte(juanMarshaled))
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			FetchClient: server.Client(),
			Address:     server.URL[7:],
		}}
		fetchParams := RecordingsFetchParams{}

		recordings, err := rm.Fetch(fetchParams)

		require.NoError(t, err)
		require.Equal(t, 2, len(recordings))
	})
}

func Test_RecordingsClient_Download(t *testing.T) {
	t.Run("returns error when max retry reached because of no recordings found", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/flv.cgi`, req.URL.Path)
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &RecordingsClient{&client{
			DownloadClient: server.Client(),
			Address:        server.URL[7:],
		}}

		var dst bytes.Buffer
		recMeta := RecordingMeta{}

		err := rm.Download(recMeta, &dst)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})
}
