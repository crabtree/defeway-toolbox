package defewayclient

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SnapshotClient_Fetch(t *testing.T) {
	t.Run("downloads the snapshot successfuly", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			require.Equal(t, `/cgi-bin/snapshot.cgi`, req.URL.Path)
			rw.Write([]byte("Hello!"))
		}))
		defer server.Close()

		rm := &SnapshotClient{
			client: fixClient(server.Client(), server.URL[7:]),
		}

		var dst bytes.Buffer
		err := rm.Fetch(0, &dst)

		require.NoError(t, err)
		require.Equal(t, "Hello!", dst.String())
	})

	t.Run("returns error when error occures during http call", func(t *testing.T) {
		server := httptest.NewServer(
			http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Write([]byte{})
			}))
		defer server.Close()

		rm := &SnapshotClient{
			client: fixClient(server.Client(), "invalid-address"),
		}

		var dst bytes.Buffer
		err := rm.Fetch(0, &dst)

		require.Contains(t, err.Error(), "dial tcp: lookup invalid-address")
	})
}
