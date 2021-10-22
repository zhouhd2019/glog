package glog

import (
	"fmt"
	"testing"
	"time"
)

func TestSimpleLog(t *testing.T) {
	cancelLogFunc, err := InitLoggerSystem("./", "test_log")
	if err != nil {
		t.Error("InitLoggerSystem Failed")
		fmt.Println(err.Error())
		return
	}
	go func() {
		logger := GetSimpleLogger("test")
		for i := 0; i < 10000; i++ {
			logger.Debug("ddd")
			logger.Info("iii")
			logger.Warning("www")
			logger.Error("eee")
			logger.Critical("ccc")
		}
	}()

	logger := GetSimpleLogger("test2")
	for i := 0; i < 10000; i++ {
		logger.Debug("1")
		logger.Info("2")
		logger.Warning("3")
		logger.Error("4")
		logger.Critical("5")
	}

	time.Sleep(10 * time.Millisecond)
	cancelLogFunc()
	time.Sleep(10 * time.Millisecond)
	CloseLoggerSystem()
}
