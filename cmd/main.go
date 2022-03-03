package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tools-home/internal/di"

	"github.com/sirupsen/logrus"
)

// @title tools-home
// @version 1.0.0
// @description tools-home created by bedrock
// @license.name MIT
func main() {
	flag.Parse()
	logrus.Info("tools-home start")

	_, closeFunc, err := di.InitApp()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		logrus.Infof("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeFunc()
			logrus.Info("tools-home exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
