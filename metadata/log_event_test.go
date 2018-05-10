package metadata

import (
	"testing"
)

// NOTE: These tests are not the best they could be since they expect the JSON to always
// be built the same way. Since JSON is order-independent, that's a bad expectation. This
// works for now, but a better solution is needed in the future.

// The following tests building a log event with system context; the intent of this test is to make sure that the
// JSON payload does not carry the "empty" child attributes of event.Context (event.Context.Platform and
// event.Context.Source).
func TestLogEventEncodeJSONEmptyChild(test *testing.T) {
	expected := `{"$schema":"https://raw.githubusercontent.com/timberio/log-event-json-schema/v3.0.8/schema.json","context":{"system":{"hostname":"localhost"}}}`
	event := NewLogEvent()
	event.ensureSystemContext()
	event.Context.System.Hostname = "localhost"
	encodedBytes, err := event.EncodeJSON()

	if err != nil {
		test.Fatal("Could not encode the event as JSON!")
	}

	encodedString := string(encodedBytes)

	if encodedString != expected {
		test.Fatalf("Expected %s but got %s", expected, encodedString)
	}
}
