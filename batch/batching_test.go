package batch

import (
	"bytes"
	"testing"
	"time"
)

func TestChannelClosing(t *testing.T) {
	lines := make(chan string)

	batcher := NewBatcher(lines, Config{
		Period: 10 * time.Second,
	})

	batcher.Lines <- "test log line"
	close(batcher.Lines)

	actual := <-batcher.BufferChan
	expected := "test log line\n"
	if actual.String() != expected {
		t.Fatalf("expected \"%+v\", got \"%+v\"", expected, actual)
	}
}

func TestBufferOverflow(t *testing.T) {
	lines := make(chan string)
	batcher := NewBatcher(lines, DefaultConfig())

	filler := "test log line"
	fillerLen := len(filler) + 1
	for written := 0; written+fillerLen < 990000; written += fillerLen {
		batcher.Lines <- filler
	}
	batcher.Lines <- "overflowed"
	close(batcher.Lines)

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
	lines := make(chan string)
	batcher := NewBatcher(lines, DefaultConfig())

	filler := "test log line"
	buf := bytes.NewBuffer(make([]byte, defaultBufferSize))
	for buf.Len() < defaultBufferSize {
		buf.WriteString(filler)
	}
	logline := buf.String()

	// go Batch(lines, bufChan, 10)
	batcher.Lines <- logline
	close(batcher.Lines)

	// Nothing should be sent to bufChan since we are dropping message
	actual := <-batcher.BufferChan
	if actual != nil {
		t.Fatalf("expected \"%+v\" to be nil", actual)
	}
}
