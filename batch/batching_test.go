package batch

import (
	"bytes"
	"testing"
	"time"
)

func TestChannelClosing(t *testing.T) {
	byteChan := make(chan []byte)

	batcher := NewBatcher(byteChan, Config{
		Period: 10 * time.Second,
	})

	batcher.ByteChan <- []byte("test log line")
	close(batcher.ByteChan)

	actual := <-batcher.BufferChan
	expected := "test log line\n"
	if actual.String() != expected {
		t.Fatalf("expected \"%+v\", got \"%+v\"", expected, actual)
	}
}

func TestBufferOverflow(t *testing.T) {
	byteChan := make(chan []byte)
	batcher := NewBatcher(byteChan, DefaultConfig())

	filler := []byte("test log line")
	fillerLen := len(filler) + 1
	for written := 0; written+fillerLen < 990000; written += fillerLen {
		batcher.ByteChan <- filler
	}
	batcher.ByteChan <- []byte("overflowed")
	close(batcher.ByteChan)

	<-batcher.BufferChan // throw away the big one
	actual := <-batcher.BufferChan
	expected := "overflowed\n"
	if actual.String() != expected {
		t.Fatalf("expected \"%+v\", got \"%+v\"", expected, actual)
	}
}

// Batch()
// Log lines larger than the max payload size (1 MB) should be dropped
func TestBatchDropLogLine(t *testing.T) {
	byteChan := make(chan []byte)
	batcher := NewBatcher(byteChan, DefaultConfig())

	filler := "test log line"
	buf := bytes.NewBuffer(make([]byte, defaultBufferSize))
	for buf.Len() < defaultBufferSize {
		buf.WriteString(filler)
	}
	logline := buf.String()

	// go Batch(lines, bufChan, 10)
	batcher.ByteChan <- []byte(logline)
	close(batcher.ByteChan)

	// Nothing should be sent to bufChan since we are dropping message
	actual := <-batcher.BufferChan
	if actual != nil {
		t.Fatalf("expected \"%+v\" to be nil", actual)
	}
}
