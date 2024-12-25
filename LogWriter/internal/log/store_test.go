package log

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	message = []byte("Hello World")
	width   = uint64(len(message)) + 8
)

func TestStore(t *testing.T) {
	f, err := os.CreateTemp("", "store_append_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	s, err := NewStore(f)
	require.NoError(t, err)

	var currentPos uint64 = 0
	//Test Append
	n, pos, err := s.Append(message)
	require.NoError(t, err)
	require.Equal(t, width, n)
	require.Equal(t, pos, currentPos)

	//Test Read
	read, err := s.Read(currentPos)
	require.NoError(t, err)
	require.Equal(t, message, read)

	//Test ReadAt
	//First reading the size of message
	var offset int64 = 0
	readSize := make([]byte, lenWidth)
	numBytesRead, err := s.ReadAt(readSize, offset)
	require.NoError(t, err)
	require.Equal(t, lenWidth, numBytesRead)
	require.Equal(t, uint64(len(message)), enc.Uint64(readSize))

	//Reading the actual message
	offset += int64(numBytesRead)
	readMessage := make([]byte, enc.Uint64(readSize))
	numBytesRead, err = s.ReadAt(readMessage, offset)
	require.NoError(t, err)
	require.Equal(t, len(message), numBytesRead)
	require.Equal(t, message, readMessage)
}
