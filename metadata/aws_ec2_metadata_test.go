package metadata

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Available()
// When the endpoint is available, Available() should return `true`
func TestEC2ClientAvailableTrue(test *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	ec2Client := NewEC2Client(DefaultConfig())
	ec2Client.BaseEndpoint = ts.URL

	available := ec2Client.Available()

	if available != true {
		test.Fatal("Expected connection to metadata provider to succeed")
	}
}

// Available()
// When the endpoint is available, but returns 404, Available() should return `false`
func TestEC2ClientAvailable404(test *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))

	ec2Client := NewEC2Client(DefaultConfig())
	ec2Client.BaseEndpoint = ts.URL

	available := ec2Client.Available()

	if available != false {
		test.Fatal("Expected connection to metadata provider to fail")
	}
}

// Available()
// When the endpoint is not available, the connection should timeout
// and Available() should return `false`
func TestEC2ClientAvailableFalse(test *testing.T) {
	ec2Client := NewEC2Client(DefaultConfig())

	available := ec2Client.Available()

	if available != false {
		test.Fatal("Expected connection to metadata provider to fail")
	}
}

// GetMetadata()
// When the service is available, the metadata should be fetched and returned
// Tests that the client hits the appropriate endpoint
func TestEC2ClientGetMetadata(test *testing.T) {
	expected := "i1934195190"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedURI := "/latest/meta-data/instance-id"
		if r.RequestURI != expectedURI {
			test.Fatalf("Expected request URI to be %s, but got %s", expectedURI, r.RequestURI)
		}

		w.WriteHeader(200)
		w.Write([]byte(expected))
	}))

	ec2Client := NewEC2Client(DefaultConfig())
	ec2Client.BaseEndpoint = ts.URL

	instanceID, err := ec2Client.GetMetadata("instance-id")

	if err != nil {
		test.Fatalf("Expected to get metadata, encountered error instead: %s", err)
	}

	if instanceID != expected {
		test.Fatalf("Expected instance ID of %s, got %s instead", expected, instanceID)
	}
}

// GetMetadata()
// When the service is available, the metadata should be fetched and returned
// Tests that the client properly handles a 404 from the service
func TestEC2ClientGetMetadata404(test *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(""))
	}))

	ec2Client := NewEC2Client(DefaultConfig())
	ec2Client.BaseEndpoint = ts.URL

	_, err := ec2Client.GetMetadata("instance-id")

	if err == nil {
		test.Fatalf("Expected to get an error when fetching metadata but didn't")
	}
}

// AddEC2Metadata()
// When the service is not available, the LogEvent should not be modified
func TestAddEC2MetadataNoService(test *testing.T) {
	ec2Client := NewEC2Client(DefaultConfig())
	logEvent := &LogEvent{}

	AddEC2Metadata(ec2Client, logEvent)

	if logEvent.Context != nil {
		test.Fatal("Expected logEvent.Context to be nil but it was not")
	}
}

// AddEC2Metadata()
// When the service is available, modifies the LogEvent with the appropriate data
func TestAddEC2Metadata(test *testing.T) {
	expected := "i1934195190"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(expected))
	}))

	ec2Client := NewEC2Client(DefaultConfig())
	ec2Client.BaseEndpoint = ts.URL
	logEvent := &LogEvent{}

	AddEC2Metadata(ec2Client, logEvent)

	instanceID := logEvent.Context.Platform.AWSEC2.InstanceID

	if instanceID != expected {
		test.Fatalf("Expected InstanceID to be %s, instead got %s", expected, instanceID)
	}
}
