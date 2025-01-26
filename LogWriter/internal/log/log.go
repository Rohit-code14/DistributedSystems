package log

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/Rohit-code14/LogWriter/internal/api"
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
	return l, l.setup()
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
		if err := l.newSegment(baseOffsets[i]); err != nil {
			return err
		}
		i++
	}

	if l.segments == nil {
		if err := l.newSegment(l.Config.Segment.InitialOffset); err != nil {
			return err
		}
	}
	return nil
}

func (l *Log) newSegment(offset uint64) error {
	s, err := newSegment(l.Dir, offset, l.Config)
	if err != nil {
		return err
	}
	l.segments = append(l.segments, s)
	l.activeSegment = s
	return nil
}

func (l *Log) Append(record *api.Record) (uint64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	off, err := l.activeSegment.Write(record)
	if err != nil {
		return 0, err
	}
	if l.activeSegment.IsMaxed() {
		err = l.newSegment(off + 1)
	}
	return off, err
}

func (l *Log) Read(offset uint64) (*api.Record, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	var s *segment
	for _, segment := range l.segments {
		if segment.baseOffset <= offset && offset < segment.nextOffset {
			s = segment
			break
		}
	}
	if s == nil || s.nextOffset <= offset {
		return nil, fmt.Errorf("Offset(%d) out of range", offset)
	}
	return s.Read(offset)
}

func (l *Log) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, segment := range l.segments {
		if err := segment.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Log) Remove() error {
	if err := l.Close(); err != nil {
		return err
	}
	return os.RemoveAll(l.Dir)
}

func (l *Log) Reset() error {
	if err := l.Remove(); err != nil {
		return err
	}
	return l.setup()
}

func (l *Log) LowestOffset() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.segments[0].baseOffset
}

func (l *Log) HighestOffset() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	offset := l.segments[len(l.segments)-1].nextOffset
	if offset == 0 {
		return 0
	}
	return offset - 1
}
