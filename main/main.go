package main

import (
	"log"
	"time"

	"github.com/hpgood/applog"
)

func main() {

	log.Println("start test log")
	applog.Info("test","hello",-1)
	time.Sleep(time.Second*5)
	applog.Info("test","world",-1)
	time.Sleep(time.Second*5)
	applog.Info("test","finish!",-1)
}