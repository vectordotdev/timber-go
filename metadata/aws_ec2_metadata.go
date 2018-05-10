package metadata

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/timberio/timber-go/logging"
)

var (
	defaultTimeout  = 1 * time.Second
	defaultEndpoint = "http://169.254.169.254"
)

type Config struct {
	Timeout time.Duration

	Logger logging.Logger
}

type EC2Client struct {
	BaseEndpoint string
	HTTPClient   *http.Client

	Config
}

func DefaultConfig() Config {
	return Config{
		Timeout: defaultTimeout,
		Logger:  logging.DefaultLogger,
	}
}

func NewEC2Client(config Config) *EC2Client {
	defaultConfig := DefaultConfig()

	if config.Timeout == 0 {
		config.Timeout = defaultConfig.Timeout
	}

	if config.Logger == nil {
		config.Logger = defaultConfig.Logger
	}

	client := &http.Client{
		Timeout: config.Timeout,
	}

	return &EC2Client{
		BaseEndpoint: defaultEndpoint,
		HTTPClient:   client,
		Config:       config,
	}
}

func AddEC2Metadata(client *EC2Client, logEvent *LogEvent) {
	if !client.Available() {
		client.Logger.Print("Agent is not running on an EC2 instance")
		return
	} else {
		client.Logger.Print("Agent is running on an EC2 instance")
	}

	amiID, err := client.GetMetadata("ami-id")

	if err != nil {
		client.Logger.Print("Could not determine AMI ID the EC2 instance was launched with")
	} else {
		client.Logger.Printf("Discovered AMI ID from AWS EC2 metadata: %s", amiID)
	}

	hostname, err := client.GetMetadata("hostname")

	if err != nil {
		client.Logger.Print("Cloud not determine the AWS assigned hostname for the EC2 instance")
	} else {
		client.Logger.Printf("Discovered EC2 Hostname from AWS EC2 metadata: %s", hostname)
	}

	instanceID, err := client.GetMetadata("instance-id")

	if err != nil {
		client.Logger.Print("Could not determine the instance ID for the EC2 instance")
	} else {
		client.Logger.Printf("Discovered Instance ID from AWS EC2 metadata: %s", instanceID)
	}

	instanceType, err := client.GetMetadata("instance-type")

	if err != nil {
		client.Logger.Print("Could not determine the instance type for the EC2 instance")
	} else {
		client.Logger.Printf("Discovered Instance Type from AWS EC2 metadata: %s", instanceType)
	}

	publicHostname, err := client.GetMetadata("public-hostname")

	if err != nil {
		client.Logger.Print("Could not determine the AWS assigned public hostname for the EC2 instance")
	} else {
		client.Logger.Printf("Discovered EC2 Public Hostname from AWS EC2 metadata: %s", publicHostname)
	}

	context := &AWSEC2Context{
		AmiID:          amiID,
		Hostname:       hostname,
		InstanceID:     instanceID,
		InstanceType:   instanceType,
		PublicHostname: publicHostname,
	}

	logEvent.AddEC2Context(context)
}

func (client *EC2Client) Available() bool {
	resp, err := client.HTTPClient.Get(client.BaseEndpoint + "/latest/meta-data/")

	if err != nil {
		return false
	}

	resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false
	}

	return true
}

func (client *EC2Client) GetMetadata(field string) (string, error) {
	resp, err := client.HTTPClient.Get(client.BaseEndpoint + "/latest/meta-data/" + field)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("Did not received a valid response for EC2 metadata")
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	return string(body[:]), nil
}
