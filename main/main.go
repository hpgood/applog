package main

import (
	"log"
	"time"

	"github.com/hpgood/applog"
)

func main() {
  log.Println("start test log")
  var userID int64=1
  applog.Fine("tag1","hello message",userID)
  applog.Info("tag1","hello message",userID)
  time.Sleep(time.Second*5)
  applog.Warn("tag2","my warn message",userID)
  applog.Error("tag2","my error message",userID)
  time.Sleep(time.Second*5)
  applog.Info("tag3","finish!",userID)
}