package ext

import (
	"bufio"
	"os"
)

type bufWriter struct {
	filename string
	file     *os.File
	bufw     *bufio.Writer
}

func (b *bufWriter) Open(filename string) error {
	var err error
	b.filename = filename
	b.file, err = os.Create(filename)
	if err != nil {
		return err
	}
	b.bufw = bufio.NewWriter(b.file)
	return nil
}

func (b *bufWriter) Close() error {
	if err := b.bufw.Flush(); err != nil {
		return err
	}
	if err := b.file.Close(); err != nil {
		return err
	}
	return nil
}

func (b *bufWriter) Write(d []byte) (int, error) {
	return b.bufw.Write(d)
}

func NewBufWriter(filename string) (*bufWriter, error) {
	b := &bufWriter{}
	if err := b.Open(filename); err != nil {
		return nil, err
	}
	return b, nil
}
