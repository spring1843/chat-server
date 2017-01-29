package chat

import (
	"io"
	"log"
	"time"
)

// SetLogFile a log file for this server and makes this server able to log
// Use server.Log() to send logs to this file
func (s *Server) SetLogFile(file io.Writer) {
	logger := new(log.Logger)
	logger.SetOutput(file)
	s.Logger = logger
	s.CanLog = true
}

// LogPrintf is a centralized logging function, so that all logs go to the same file and they all have time stamps
// Ads a time stamp to every log entry
// For readability start the message with a category followed by \t
func (s *Server) LogPrintf(format string, v ...interface{}) {
	if s.CanLog != true {
		return
	}
	now := time.Now()
	s.Logger.Printf(now.Format(time.UnixDate)+"\t"+format, v...)
}
