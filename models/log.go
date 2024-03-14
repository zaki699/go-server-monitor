package models

// Log ...
type Log struct {
	Level     string
	Message   string
	Module    string
	SessionId int
}

func NewLog() Log {

	return Log{
		Level:     "unknown",
		Module:    "unknown",
		SessionId: -1,
	}
}

// GetLevel ...
func (l Log) GetLevel() string {
	return l.Level
}

// GetMessage ...
func (l Log) GetMessage() string {
	return l.Message
}

// GetModule ...
func (l Log) GetModule() string {
	return l.Module
}

// GetSessionId ...
func (l Log) GetSessionId() int {
	return l.SessionId
}
