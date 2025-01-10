package log

import (
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var offsetWidth uint64 = 4
var positionWidth uint64 = 8
var entireWidth uint64 = offsetWidth + positionWidth

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}
	idx.size = uint64(fi.Size())
	if err = os.Truncate(
		f.Name(), int64(c.Segment.MaxIndexBytes),
	); err != nil {
		return nil, err
	}
	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}
	return idx, nil
}

func (i *index) Name() string {
	return i.file.Name()
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}
	if err := i.file.Sync(); err != nil {
		return err
	}
	if err := i.file.Truncate(int64(i.size)); err != nil {
		return err
	}
	return i.file.Close()
}

func (i *index) Read(offset int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}
	if offset == -1 {
		out = uint32((i.size / entireWidth) - 1)
	} else {
		out = uint32(offset)
	}
	pos = uint64(out) * entireWidth
	if i.size < pos+entireWidth {
		return 0, 0, io.EOF
	}
	out = enc.Uint32(i.mmap[pos : pos+offsetWidth])
	pos = enc.Uint64(i.mmap[pos+offsetWidth : pos+entireWidth])

	return out, pos, nil
}

func (i *index) Write(offset uint32, pos uint64) error {
	if uint64(len(i.mmap)) < i.size+entireWidth {
		return io.EOF
	}
	enc.PutUint32(i.mmap[i.size:i.size+offsetWidth], offset)
	enc.PutUint64(i.mmap[i.size+offsetWidth:i.size+entireWidth], pos)
	i.size += uint64(entireWidth)
	return nil
}
