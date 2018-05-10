package metadata

import (
	"encoding/json"
)

var schema string = "https://raw.githubusercontent.com/timberio/log-event-json-schema/v3.0.8/schema.json"

type LogEvent struct {
	Schema  string   `json:"$schema"`
	Context *Context `json:"context,omitempty"`
}

type Context struct {
	System   *SystemContext   `json:"system,omitempty"`
	Platform *PlatformContext `json:"platform,omitempty"`
	Source   *SourceContext   `json:"source,omitempty"`
}

type SystemContext struct {
	Hostname string `json:"hostname,omitempty"`
}

type PlatformContext struct {
	AWSEC2 *AWSEC2Context `json:"aws_ec2,omitempty"`
}

type SourceContext struct {
	FileName string `json:"file_name,omitempty"`
}

type AWSEC2Context struct {
	AmiID          string `json:"ami_id,omitempty"`
	Hostname       string `json:"hostname,omitempty"`
	InstanceID     string `json:"instance_id,omitempty"`
	InstanceType   string `json:"instance_type,omitempty"`
	PublicHostname string `json:"public_hostname,omitempty"`
}

func NewLogEvent() *LogEvent {
	return &LogEvent{Schema: schema}
}

func (logEvent *LogEvent) AddEC2Context(context *AWSEC2Context) {
	logEvent.ensurePlatformContext()
	logEvent.Context.Platform.AWSEC2 = context
}

func (logEvent *LogEvent) AddHostname(hostname string) {
	logEvent.ensureSystemContext()
	logEvent.Context.System.Hostname = hostname
}

func (logEvent *LogEvent) EncodeJSON() ([]byte, error) {
	return json.Marshal(logEvent)
}

func (logEvent *LogEvent) ensureContext() {
	if logEvent.Context == nil {
		logEvent.Context = &Context{}
	}
}

func (logEvent *LogEvent) ensurePlatformContext() {
	logEvent.ensureContext()
	if logEvent.Context.Platform == nil {
		logEvent.Context.Platform = &PlatformContext{}
	}
}

func (logEvent *LogEvent) ensureSystemContext() {
	logEvent.ensureContext()
	if logEvent.Context.System == nil {
		logEvent.Context.System = &SystemContext{}
	}
}

func (logEvent *LogEvent) ensureSourceContext() {
	logEvent.ensureContext()
	if logEvent.Context.Source == nil {
		logEvent.Context.Source = &SourceContext{}
	}
}
