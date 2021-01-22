# applog
applog to cloud using http.

## 包名：
"github.com/hpgood/applog"

配置文件data/config/app.ini
加入：
```
[log]
project=applog_demo
version=0.0.1
time=5000
appid=100
token=from_k8s
url=http://apps-applog.default/logcat/server
```

## 说明：
```
project 工程名字
version 版本号
time 日志提交频率,默认 5000 ms
appid 自定义整数appid 建议用3位固定数字,减少重复。
token from_k8s 默认即可
url  http://apps-applog.default/logcat/server  默认即可
```

## 代码例子：

https://github.com/hpgood/applog/blob/main/main/main.go
```
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
```

# JS
```
<script src="js/applog.min.js"></script>
<script>

// 888 是appID
// my_project_name 工程名字
//
applog.init(888,"my_project_name","js/applog-worker.min.js");
// console 重定向默认false
applog.setRedirect(false);

//设置用户ID
applog.setUserID(123);

// info
applog.info("main.html","hello info message");
// warn
applog.warn("main.html","warn message");
// error
applog.error("err.html","error message");

</script>
```