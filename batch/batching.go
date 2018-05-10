package batch

import (
	"bytes"
	"io"
	"log"
	"os"
	"time"
)

var (
	// Preallocate 990kb. The Timber API will not accept payloads larger than 1mb.
	// This leaves 10kb for headers.
	defaultBufferSize = 990000

	defaultPeriod = 3 * time.Second
)

type Batcher struct {
	BufferChan chan *bytes.Buffer
	Lines      chan string

	Config
}

type Config struct {
	Period time.Duration
	Size   int

	Logger *log.Logger
}

func DefaultConfig() Config {
	return Config{
		Period: defaultPeriod,
		Size:   defaultBufferSize,

		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func NewBatcher(lines chan string, config Config) *Batcher {
	defaultConfig := DefaultConfig()

	if config.Period == 0 {
		config.Period = defaultConfig.Period
	}

	if config.Size == 0 {
		config.Size = defaultConfig.Size
	}

	if config.Logger == nil {
		config.Logger = defaultConfig.Logger
	}

	batcher := &Batcher{
		BufferChan: make(chan *bytes.Buffer),
		Lines:      lines,
		Config:     config,
	}

	go batcher.batch()

	return batcher
}

func Batch(lines chan string) *Batcher {
	batcher := &Batcher{
		BufferChan: make(chan *bytes.Buffer),
		Lines:      lines,
		Config:     DefaultConfig(),
	}

	go batcher.batch()

	return batcher
}

func (batcher *Batcher) batch() {
	buffer := freshBuffer(batcher.Size)
	tick := time.Tick(batcher.Period)

	for {
		select {
		case line, ok := <-batcher.Lines:
			if ok {
				if len(line)+1 > buffer.Cap() {
					// @TODO track souce?
					// @TODO more informative log line
					// have callback channel for results
					batcher.Logger.Printf("Dropping log line greater than the max buffer size")
					continue
				}

				if buffer.Len()+len(line)+1 > buffer.Cap() {
					batcher.BufferChan <- buffer
					buffer = freshBuffer(batcher.Size)
				}

				if len(line) > 0 {
					io.WriteString(buffer, line+"\n")
				}

			} else { // channel is closed
				if buffer.Len() > 0 {
					batcher.BufferChan <- buffer
				}
				close(batcher.BufferChan)
				return
			}

		case <-tick:
			if buffer.Len() > 0 {
				batcher.BufferChan <- buffer
				buffer = freshBuffer(batcher.Size)
			}
		}
	}
}

func freshBuffer(size int) *bytes.Buffer {
	buf := bytes.NewBuffer(make([]byte, size))
	buf.Reset()
	return buf
}
