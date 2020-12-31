# applog
applog to cloud using http.

## 包名：
"github.com/hpgood/applog"

#配置文件data/config/app.ini
加入：

[log]
project=applog_demo
version=0.0.1
time=5000
appid=100
token=from_k8s
url=http://apps-applog.default/logcat/server

## 说明：
project 工程名字
version 版本
time 日志提交频率,默认 5000 ms
appid 自定义整数appid 建议用3位固定数字,减少重复。
token from_k8s 默认即可
url  http://apps-applog.default/logcat/server  默认即可，调试可以用 https://api.yondor.cn/logcat/server

## 代码例子：

https://github.com/hpgood/applog/blob/main/main/main.go

import (
	"log"
	"time"

	"github.com/hpgood/applog"
)

func main() {

	log.Println("start test log")
  //Fine 级别日志
  applog.Fine("test","hello",-1)
  //Info 级别日志
	applog.Info("test","hello",-1)
	time.Sleep(time.Second*5)
  //Warn 级别日志
	applog.Warn("test","world",-1)
  //Error 级别日志
  applog.Error("test","world",-1)
	time.Sleep(time.Second*5)
	applog.Info("test","finish!",-1)
}

