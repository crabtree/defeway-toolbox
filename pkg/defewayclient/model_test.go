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
}

func TestUnmarshalJuanForRecSearch(t *testing.T) {
	t.Run("unmarshal XML from bytes with empty RecSearch", func(t *testing.T) {
		juanMarshaled := `<juan ver="" squ="" dir="0" enc="0" errno="0"></juan>`
		juan, err := UnmarshalJuanForRecSearch([]byte(juanMarshaled))

		require.NoError(t, err)
		require.Nil(t, juan.RecSearch)
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
		juan, err := UnmarshalJuanForRecSearch([]byte(juanMarshaled))

		require.NoError(t, err)
		require.NotNil(t, juan.RecSearch)
		validateRecSearch(t, juan.RecSearch)
		validateJuan(t, juan)
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
