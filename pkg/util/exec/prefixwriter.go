package exec

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

type PrefixWriter struct {
	pr   *io.PipeReader
	pw   *io.PipeWriter
	done chan struct{}
}

func NewPrefixWriter(w io.Writer, prefix string) *PrefixWriter {
	pw := &PrefixWriter{
		done: make(chan struct{}),
	}
	pw.pr, pw.pw = io.Pipe()

	go func() {
		defer close(pw.done)
		s := bufio.NewScanner(pw.pr)
		for s.Scan() {
			line := fmt.Sprintf("%s%s\n", prefix, s.Bytes())
			//nolint:errcheck
			w.Write([]byte(line))
		}
	}()
	return pw
}

func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	return pw.pw.Write(p)
}

func (pw *PrefixWriter) Close() error {
	pw.pw.Close()
	pw.pr.Close()
	select {
	case <-time.After(time.Second):
	case <-pw.done:
	}
	return nil
}
