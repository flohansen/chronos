package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/flohansen/chronos/internal/metric"
)

type Writer struct {
	f io.Writer
}

func NewFileWriter(f io.Writer) *Writer {
	return &Writer{
		f: f,
	}
}

func (w *Writer) Write(m metric.Metric) error {
	timestamp := time.Now().UnixMilli()
	row := fmt.Sprintf("%s@%d=%.2f\n", m.Name, timestamp, m.Value)

	if _, err := w.f.Write([]byte(row)); err != nil {
		return fmt.Errorf("could not write: %s", err)
	}

	return nil
}
