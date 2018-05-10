package forward

import (
	"bytes"
)

type Forwarder interface {
	Forward(*bytes.Buffer)
}

func Forward(bufferChan chan *bytes.Buffer, forwarder Forwarder) {
	for buffer := range bufferChan {
		forwarder.Forward(buffer)
	}
}
