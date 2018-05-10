package forward

import (
	"bufio"
	"bytes"
	"os"

	"github.com/timberio/timber-go/logging"
)

type StdoutForwarder struct {
	Logger logging.Logger
}

func NewStdoutForwarder(logger logging.Logger) *StdoutForwarder {
	if logger == nil {
		logger = logging.DefaultLogger
	}

	return &StdoutForwarder{
		Logger: logger,
	}
}

func (s *StdoutForwarder) Forward(buffer *bytes.Buffer) {
	writer := bufio.NewWriter(os.Stdout)

	_, err := writer.Write(buffer.Bytes())
	if err != nil {
		s.Logger.Print(err)
		return
	}

	err = writer.Flush()
	if err != nil {
		s.Logger.Print(err)
		return
	}
}
