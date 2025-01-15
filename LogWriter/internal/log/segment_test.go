package log

import (
	"io"
	"os"
	"testing"

	"github.com/Rohit-code14/LogWriter/internal/api"
	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {
	dir, err := os.MkdirTemp("", "segment_test")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	c.Segment.MaxStoreBytes = entireWidth * 3

	s, err := newSegment(dir, 16, c)
	require.NoError(t, err)
	require.Equal(t, uint64(16), s.baseOffset, s.nextOffset)
	require.False(t, s.IsMaxed())

	record := &api.Record{Value: []byte("Logging")}

	for i := uint64(0); i < 3; i++ {
		offset, err := s.Write(record)

		require.NoError(t, err)
		require.Equal(t, 16+i, offset)

		message, err := s.Read(offset)
		require.NoError(t, err)
		require.Equal(t, record.Value, message.Value)
	}

	_, err = s.Write(record)
	require.Error(t, io.EOF, err)
	require.True(t, s.IsMaxed())

	c.Segment.MaxStoreBytes = uint64(len(record.Value) * 3)
	c.Segment.MaxIndexBytes = 1024

	s, err = newSegment(dir, 16, c)
	require.NoError(t, err)
	require.True(t, s.IsMaxed())
	err = s.Remove()
	require.NoError(t, err)

}
