package log

import (
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Log struct {
	mu            sync.RWMutex
	Dir           string
	Config        Config
	activeSegment *segment
	segments      []*segment
}

func NewLog(dir string, config Config) (*Log, error) {
	if config.Segment.MaxIndexBytes == 0 || config.Segment.MaxStoreBytes == 0 {
		config.Segment.MaxIndexBytes = 1024
		config.Segment.MaxStoreBytes = 1024
	}
	l := &Log{
		Dir:    dir,
		Config: config,
	}
	return l, nil
}

func (l *Log) setup() error {
	files, err := os.ReadDir(l.Dir)
	if err != nil {
		return err
	}
	var baseOffsets []uint64
	for _, file := range files {
		offStr := strings.TrimSuffix(file.Name(), path.Ext(file.Name()))
		offSet, _ := strconv.ParseUint(offStr, 10, 0)
		baseOffsets = append(baseOffsets, offSet)
	}

	sort.Slice(baseOffsets, func(i, j int) bool {
		return baseOffsets[i] < baseOffsets[j]
	})

	for i := 0; i < len(baseOffsets); i++ {

	}

	return nil
}
