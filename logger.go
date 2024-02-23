package main

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	fileName string
	file     *os.File
	mu       sync.Mutex
}

func (l *Logger) SetLogFile() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.fileName = l.generateLogFileName()
	if l.file != nil {
		l.file.Close()
	}
	l.openLogFile()
}

func (l *Logger) generateLogFileName() string {
	currentTime := time.Now()
	return filepath.Join(getLogFolderPath(), "status_"+currentTime.Format("2006_01_02_15")+".log")
}

func getLogFolderPath() string {
	logFolderPath := filepath.Join(getSavedGamesFolder(), "MFDMF", "Logs")
	return logFolderPath
}

func getBaseDirectory() string {
	return filepath.Join(getSavedGamesFolder(), "MFDMF", "Modules")
}

func getSavedGamesFolder() string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Failed to get current user: %v", err)
	}
	savedGamesFolder := filepath.Join(currentUser.HomeDir, "Saved Games")
	return savedGamesFolder
}

func (l *Logger) openLogFile() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.fileName = l.generateLogFileName()

	logFolder := filepath.Dir(l.fileName)
	err := os.MkdirAll(logFolder, 0755)
	if err != nil {
		log.Fatalf("Failed to create log folder: %v", err)
	}

	file, err := os.OpenFile(l.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	l.file = file

	log.SetOutput(l.file)
}

func (l *Logger) Log(message string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	log.Println(message)
}

var instance *Logger
var once sync.Once

func GetLogger() *Logger {
	once.Do(func() {
		instance = &Logger{}
		instance.openLogFile()
	})
	return instance
}
