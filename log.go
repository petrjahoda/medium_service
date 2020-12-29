package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func logInfo(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --INF-- " + data)
	appendDataToLog("INF", reference, data)
}

func logError(reference, data string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05.000") + " [" + reference + "] --ERR-- " + data)
	appendDataToLog("ERR", reference, data)
}

func logDirectoryCheck() {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	dir := getActualDirectory(dateTimeFormat)
	createLogDirectory(dir, dateTimeFormat)
}

func createLogDirectory(dir string, dateTimeFormat string) {
	logDirectory := filepath.Join(dir, "log")
	_, checkPathError := os.Stat(logDirectory)
	logDirectoryExists := checkPathError == nil
	if logDirectoryExists {
		return
	}
	switch runtime.GOOS {
	case "windows":
		{
			err := os.Mkdir(logDirectory, 0777)
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to create directory for log file: " + err.Error())
				return
			}
			fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --INF-- Log directory created")
		}

	default:
		{
			err := os.MkdirAll(logDirectory, 0777)
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to create directory for log file: " + err.Error())
				return
			}
			fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --INF-- Log directory created")
		}
	}
}

func getActualDirectory(dateTimeFormat string) string {
	var dir string
	switch runtime.GOOS {
	case "windows":
		{
			executable, err := os.Executable()
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to read actual directory: " + err.Error())
			}
			dir = filepath.Dir(executable)
		}
	default:
		{
			executable, err := os.Getwd()
			if err != nil {
				fmt.Println(time.Now().Format(dateTimeFormat) + " [MAIN] --ERR-- Unable to read actual directory: " + err.Error())
			}
			dir = executable
		}
	}
	return dir
}

func appendDataToLog(logLevel string, reference string, data string) {
	dateTimeFormat := "2006-01-02 15:04:05.000"
	logNameDateTimeFormat := "2006-01-02"
	logDirectory := filepath.Join(".", "log")
	logFileName := reference + " " + time.Now().Format(logNameDateTimeFormat) + ".log"
	logFullPath := strings.Join([]string{logDirectory, logFileName}, "/")
	logData := time.Now().Format("2006-01-02 15:04:05.000 ") + reference + " " + logLevel + " " + data
	writingSync.Lock()
	f, err := os.OpenFile(logFullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot open file: " + err.Error())
		writingSync.Unlock()
		return
	}
	_, err = f.WriteString(logData + "\r\n")
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot write to file: " + err.Error())
	}
	err = f.Close()
	if err != nil {
		fmt.Println(time.Now().Format(dateTimeFormat) + " [" + reference + "] --ERR-- Cannot close file: " + err.Error())
	}
	writingSync.Unlock()
}

func deleteOldLogFiles(deleteLogsAfter time.Duration) {
	for serviceIsRunning {
		directory, err := ioutil.ReadDir("log")
		if err != nil {
			logError("MAIN", "Problem opening log directory")
			return
		}
		now := time.Now()
		logDirectory := filepath.Join(".", "log")
		for _, file := range directory {
			if fileAge := now.Sub(file.ModTime()); fileAge > deleteLogsAfter {
				logInfo("MAIN", "Deleting old log file "+file.Name()+" with age of "+fileAge.String())
				logFullPath := strings.Join([]string{logDirectory, file.Name()}, "/")
				var err = os.Remove(logFullPath)
				if err != nil {
					logError("MAIN", "Problem deleting file "+file.Name()+", "+err.Error())
				}
			}
		}
		time.Sleep(1 * time.Hour)
	}
}
