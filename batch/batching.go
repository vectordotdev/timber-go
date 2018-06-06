package batch

import (
	"bytes"
	"log"
	"os"
	"time"

	"github.com/timberio/timber-go/logging"
)

var (
	// Preallocate 990kb. The Timber API will not accept payloads larger than 1mb.
	// This leaves 10kb for headers.
	defaultBufferSize = 990000

	defaultPeriod = 3 * time.Second
)

type Batcher struct {
	BufferChan chan *bytes.Buffer
	ByteChan   chan []byte

	Config
}

type Config struct {
	Period time.Duration
	Size   int

	Logger logging.Logger
}

func DefaultConfig() Config {
	return Config{
		Period: defaultPeriod,
		Size:   defaultBufferSize,

		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func NewBatcher(byteChan chan []byte, config Config) *Batcher {
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
		ByteChan:   byteChan,
		Config:     config,
	}

	go batcher.batch()

	return batcher
}

func Batch(byteChan chan []byte) *Batcher {
	batcher := &Batcher{
		BufferChan: make(chan *bytes.Buffer),
		ByteChan:   byteChan,
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
		case b, ok := <-batcher.ByteChan:
			if ok {
				if len(b)+1 > buffer.Cap() {
					// @TODO track souce?
					// @TODO more informative log line
					// have callback channel for results
					batcher.Logger.Print("Dropping log line greater than the max buffer size")
					continue
				}

				if buffer.Len()+len(b)+1 > buffer.Cap() {
					batcher.BufferChan <- buffer
					buffer = freshBuffer(batcher.Size)
				}

				if len(b) > 0 {
					b = append(b, []byte("\n")...)
					buffer.Write(b)
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
