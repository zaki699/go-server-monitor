package models

import (
	"time"
)

// Progress ...
type Stat struct {
	StatId       int64
	Frames       int64
	DropFrames   int64
	DupFrames    int64
	Bitrate      string
	Speed        float64
	Fps          float64
	SessionId    int64
	EncodingTime int64
	StreamsQP    string
	CreatedAt    time.Time
}

func NewStat() Stat {

	return Stat{
		Frames:       -1,
		Speed:        -1.0,
		Fps:          -1.0,
		SessionId:    -1,
		DropFrames:   -1,
		DupFrames:    -1,
		EncodingTime: -1,
		CreatedAt:    time.Now(),
	}
}

// GetStatId ...
func (p Stat) GetStatId() int64 {
	return p.StatId
}

// GetCreatedAt ...
func (p Stat) GetCreatedAt() time.Time {
	return p.CreatedAt
}

// GetStreamsQP ...
func (p Stat) GetStreamsQP() string {
	return p.StreamsQP
}

// GetEncodingtime ...
func (p Stat) GetEncodingTime() int64 {
	return p.EncodingTime
}

// GetFrames ...
func (p Stat) GetFrames() int64 {
	return p.Frames
}

// GetDropFrames ...
func (p Stat) GetDropFrames() int64 {
	return p.DropFrames
}

// GetDupFrames ...
func (p Stat) GetDupFrames() int64 {
	return p.DupFrames
}

// GetBitrate ...
func (p Stat) GetBitrate() string {
	return p.Bitrate
}

// GetFPS ...
func (p Stat) GetFPS() float64 {
	return p.Fps
}

// GetSpeed ...
func (p Stat) GetSpeed() float64 {
	return p.Speed
}

// GetSessionsId ...
func (p Stat) GetSessionId() int64 {
	return p.SessionId
}
