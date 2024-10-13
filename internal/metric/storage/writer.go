package storage

import (
	"fmt"
	"os"
	"time"

	"github.com/flohansen/chronos/internal/metric"
)

type Writer struct {
}

func NewFileWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Write(m metric.Metric) error {
	filename := "current"

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("could not open file: %s", err)
	}
	defer f.Close()

	timestamp := time.Now().UnixMilli()
	row := fmt.Sprintf("%s@%d=%.2f\n", m.Name, timestamp, m.Value)

	if _, err := f.Write([]byte(row)); err != nil {
		return fmt.Errorf("could not write: %s", err)
	}

	return nil
}
