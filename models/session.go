package models

import (
	"net/url"
	"reflect"
	"time"
)

// Session ...
type Session struct {
	SessionId         int64
	ChannelName       string
	Codec             string
	Definition        string
	Preset            string
	HostName          string
	OptimizerEnabled  string
	Name              string
	Status            string
	Cmd				  string
	CreatedAt         time.Time
}

func New(url url.URL) Session {

	if reflect.ValueOf(url).IsZero() {
		panic("url cannot be nil")
	}

	return Session{
		ChannelName:       url.Query().Get("channel_name"),
		Codec:             url.Query().Get("codec"),
		Definition:        url.Query().Get("definition"),
		HostName:          url.Query().Get("hostname"),
		Preset:            url.Query().Get("preset"),
		OptimizerEnabled:  url.Query().Get("optimizer_enabled"),
		Name:              url.Query().Get("name"),
		Cmd:			   url.Query().Get("cmd"),
		Status:            "running",
		CreatedAt:         time.Now(),
	}
}

// GetCmd ...
func (s Session) GetCmd() string {
	return s.Cmd
}

// GetCreatedAt ...
func (s Session) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetSessionId ...
func (s Session) GetSessionId() int64 {
	return s.SessionId
}

// GetPreset ...
func (s Session) GetPreset() string {
	return s.Preset
}

// GetName ...
func (s Session) GetName() string {
	return s.Name
}

// GetChannelName ...
func (s Session) GetChannelName() string {
	return s.ChannelName
}

// GetCodec ...
func (s Session) GetCodec() string {
	return s.Codec
}

// GetDefinition ...
func (s Session) GetDefinition() string {
	return s.Definition
}

// GetHostName ...
func (s Session) GetHostName() string {
	return s.HostName
}

// GetOptimizerEnabled ...
func (s Session) GetOptimizerEnabled() string {
	return s.OptimizerEnabled
}

// GetStatus ...
func (s Session) GetStatus() string {
	return s.Status
}