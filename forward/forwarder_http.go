package forward

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/timberio/timber-go/logging"
)

var (
	// Added at compile time
	version string

	defaultHTTPForwarderEndpoint = "https://logs.timber.io/frames"
	defaultHTTPForwarderTimeout  = 10 * time.Second
)

type HTTPForwarder struct {
	HTTPClient *retryablehttp.Client
	APIKey     string

	Config
}

type Config struct {
	Endpoint  string
	Metadata  string
	UserAgent string

	Logger logging.Logger
}

func DefaultConfig() Config {
	return Config{
		Endpoint:  defaultHTTPForwarderEndpoint,
		UserAgent: fmt.Sprintf("timber-go-forward/%s", version),

		Logger: logging.DefaultLogger,
	}
}

func NewHTTPForwarder(apiKey string, config Config) (*HTTPForwarder, error) {
	if apiKey == "" {
		return nil, errors.New("API KEY REQUIRED")
	}

	defaultConfig := DefaultConfig()

	if config.Endpoint == "" {
		config.Endpoint = defaultConfig.Endpoint
	}

	if config.UserAgent == "" {
		config.UserAgent = defaultConfig.UserAgent
	}

	if config.Logger == nil {
		config.Logger = defaultConfig.Logger
	}

	httpClient := retryablehttp.NewClient()
	httpClient.Logger = config.Logger.(*log.Logger)
	httpClient.HTTPClient.Timeout = defaultHTTPForwarderTimeout

	return &HTTPForwarder{
		HTTPClient: httpClient,
		APIKey:     apiKey,
		Config:     config,
	}, nil
}

func (h *HTTPForwarder) Forward(buffer *bytes.Buffer) {
	token := base64.StdEncoding.EncodeToString([]byte(h.APIKey))
	authorization := fmt.Sprintf("Basic %s", token)

	req, err := retryablehttp.NewRequest("POST", h.Endpoint, bytes.NewReader(buffer.Bytes()))
	if err != nil {
		h.Logger.Print(err)
		return
	}

	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Authorization", authorization)
	req.Header.Add("User-Agent", h.UserAgent)
	req.Header.Add("Timber-Metadata-Override", h.Metadata)

	resp, err := h.HTTPClient.Do(req)
	if err != nil {
		// retries have already happened at this point, so give up
		h.Logger.Print(err)
		return
	}
	resp.Body.Close()

	if resp.StatusCode >= 300 {
		//@TODO: More informative logging message
		h.Logger.Printf("HTTPForwarder: unexpected response (status code %d)", resp.StatusCode)
		return
	}
}
