package persistence

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

type Disk struct {
	path string
	fd   *os.File
	mu   sync.Mutex
}

func (d *Disk) Persist(ctx context.Context, in int64) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	deadline, ok := ctx.Deadline()
	if ok {
		err := d.fd.SetWriteDeadline(deadline)
		if err != nil {
			fmt.Fprintf(os.Stderr, "set write deadline failed: %s\n", err)
			// ignore the deadline maybe not support by this os
		}
	}
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(in))
	n, err := d.fd.WriteAt(buf[:8], 0)
	if err != nil {
		return err
	}
	if n != 8 {
		return errors.New("incomplete clock persistent")
	}
	err = d.fd.Sync()
	if err != nil {
		return err
	}
	return nil
}

func (d *Disk) Load(ctx context.Context) (int64, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	deadline, ok := ctx.Deadline()
	if ok {
		err := d.fd.SetReadDeadline(deadline)
		if err != nil {
			fmt.Fprintf(os.Stderr, "set write deadline failed: %s\n", err)
			// ignore the deadline maybe not support by this os
		}
	}
	var buf [8]byte
	n, err := d.fd.ReadAt(buf[:8], 0)
	if n != 8 {
		// if read error but not eof, means other problem happened on read
		if !errors.Is(err, io.EOF) {
			return 0, err
		}
		return 0, ClockNotFound
	}
	return int64(binary.BigEndian.Uint64(buf[:8])), nil
}

var _ Interface = (*Disk)(nil)

func NewDisk(path string) (Interface, error) {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		return nil, err
	}
	return &Disk{
		path: path,
		fd:   fd,
	}, nil
}
