package log

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFileName string
	logMutex    sync.Mutex
)

func InitLog() {
	// Set the log file name to the current time
	SetLogFileName()

	// Start a scheduled task to update the log file every other day
	go UpdateLogFileDaily()
	// Start a scheduled task to delete historical log files every two days
	go DeleteOldLogFiles()
}

// SetLogFileName Set the log file name to the current time
func SetLogFileName() {
	logMutex.Lock()
	defer logMutex.Unlock()

	// Create a logs folder
	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		logger.Errorf("Unable to create log directory: %v", err)
	}

	// Build log file name
	logFileName = filepath.Join("logs", time.Now().Format("2006-01-02_15-04-05")+".log")

	// Set log output to file
	file, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logger.Errorf("Unable to open log file: %v", err)
	}
	logger.SetOutput(file)
}

// UpdateLogFileDaily Update the log file every other day
func UpdateLogFileDaily() {
	for {
		time.Sleep(time.Hour * 24)
		SetLogFileName()
	}
}

// DeleteOldLogFiles deletes logs older than three days every five days
func DeleteOldLogFiles() {
	for {
		time.Sleep(time.Hour * 24 * 5)
		DeleteLogsOlderThanThreeDays()
	}
}

// DeleteLogsOlderThanThreeDays deletes logs older than three days
func DeleteLogsOlderThanThreeDays() {
	logMutex.Lock()
	defer logMutex.Unlock()

	files, err := filepath.Glob("logs/*.log")
	if err != nil {
		logger.Errorf("Unable to get the list of log files: %v", err)
		return
	}

	// Get the current time minus three days
	threeDaysAgo := time.Now().Add(-time.Hour * 24 * 3)

	// Delete logs older than three days
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			logger.Errorf("Unable to get file information: %v", err)
			continue
		}

		if info.ModTime().Before(threeDaysAgo) {
			if err := os.Remove(file); err != nil {
				logger.Errorf("Unable to delete log file %s: %v", file, err)
			}
		}
	}
}
