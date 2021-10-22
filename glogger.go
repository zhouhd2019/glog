package glog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const CHECK_INTERVAL = 100
const MAX_LOG_FILE_SIZE = 1024 * 1024

var (
	logFilePath string
	logFileName string

	logFile   *os.File
	logWriter *bufio.Writer

	logChan chan *string

	logFileIdx int
	logCheckCD int
	lastDayStr string
)

type Logger interface {
	debug(logStr string)
	info(logStr string)
	warning(logStr string)
	error(logStr string)
	critical(logStr string)
}

type SimpleLogger struct {
	name string
}

//	Active log: {fileName}
//	Old Log: {fileName}.2006-01-02.{logFileIdx}
//	If file exist, clear it
//	Try to create new file on day passed Or file too large
func InitLoggerSystem(filePath string, fileName string) (context.CancelFunc, error) {
	CloseLoggerSystem()

	logFilePath = filePath
	logFileName = fileName
	logFullPath := path.Join(filePath, fileName)
	var err error
	logFile, err = os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0660))
	if err != nil {
		return nil, err
	}
	logWriter = bufio.NewWriter(logFile)
	if logWriter == nil {
		CloseLoggerSystem()
		return nil, errors.New("new writer failed")
	}
	logFileIdx = 0
	logCheckCD = 0
	lastDayStr = time.Now().Format("2006-01-02")
	logChan = make(chan *string, 1000)

	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case logStr := <-logChan:
				_, err := logWriter.Write([]byte(*logStr))
				if err != nil {
					fmt.Println(err.Error())
				}
				checkCreateNewLog()
			}

		}
	}(ctx)

	return cancel, nil
}

func ResetLoggerSystem() {
	if logWriter != nil {
		logWriter.Flush()
		logWriter = nil
	}
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

func CloseLoggerSystem() {
	if logChan != nil {
		close(logChan)
	}
	ResetLoggerSystem()
}

func GetSimpleLogger(name string) SimpleLogger {
	return SimpleLogger{name}
}

func (logger *SimpleLogger) Debug(logStr string) {
	add_log(logger.name, "[DEBUG]", logStr)
}

func (logger *SimpleLogger) Info(logStr string) {
	add_log(logger.name, "[INFO]", logStr)
}

func (logger *SimpleLogger) Warning(logStr string) {
	add_log(logger.name, "[WARNING]", logStr)
}

func (logger *SimpleLogger) Error(logStr string) {
	add_log(logger.name, "[ERROR]", logStr)
}

func (logger *SimpleLogger) Critical(logStr string) {
	add_log(logger.name, "[CRITICAL]", logStr)

}

func add_log(logStrs ...string) {
	fullLogStr := time.Now().Format("2006-01-02 15:04:05.000 ") + strings.Join(logStrs, " ") + "\n"
	logChan <- &fullLogStr
}

func checkCreateNewLog() {
	logCheckCD += 1
	if logCheckCD > CHECK_INTERVAL {
		logCheckCD = 0
		dayStr := time.Now().Format("2006-01-02")
		if lastDayStr[9] != dayStr[9] {
			createNewLog()
			lastDayStr = dayStr
			logFileIdx = 0
		} else {
			filePath := path.Join(logFilePath, logFileName)
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if fileInfo.Size() > MAX_LOG_FILE_SIZE {
				createNewLog()
				logFileIdx += 1
			}
		}
	}
}

func createNewLog() {
	handleOldLogFile()
	createNewLogFile()
}

func handleOldLogFile() {
	//	Old Log: {fileName}.2006-01-02.{logFileIdx}
	ResetLoggerSystem()
	oldPath := path.Join(logFilePath, logFileName)
	newFileNamePart := []string{logFileName, lastDayStr, strconv.Itoa(logFileIdx)}
	newFileName := strings.Join(newFileNamePart, ".")
	newPath := path.Join(logFilePath, newFileName)
	err := os.Rename(oldPath, newPath)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func createNewLogFile() {
	logFullPath := path.Join(logFilePath, logFileName)
	var err error
	logFile, err = os.OpenFile(logFullPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0660))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	logWriter = bufio.NewWriter(logFile)
	if logWriter == nil {
		ResetLoggerSystem()
		fmt.Println("new writer failed")
		return
	}
}
