package glog

import (
	"fmt"
	"testing"
	"time"

	glog "github.com/zhouhd2019/glog/glog"
)

func TestSimpleLog(t *testing.T) {
	cancelLogFunc, err := glog.InitLoggerSystem("./", "test_log")
	if err != nil {
		t.Error("InitLoggerSystem Failed")
		fmt.Println(err.Error())
		return
	}
	go func() {
		logger := glog.GetSimpleLogger("test")
		for i := 0; i < 10000; i++ {
			logger.Debug("ddd")
			logger.Info("iii")
			logger.Warning("www")
			logger.Error("eee")
			logger.Critical("ccc")
		}
	}()

	logger := glog.GetSimpleLogger("test2")
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
	glog.CloseLoggerSystem()
}
