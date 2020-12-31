# applog
applog to cloud using http.
## Golang 版本 日志收集到腾讯云

### 包名：
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

### 说明：
```
project 工程名字
version 版本
time 日志提交频率,默认 5000 ms
appid 自定义整数appid 建议用3位固定数字,减少重复。
token from_k8s 默认即可
url  http://apps-applog.default/logcat/server  默认即可
```

### 代码例子：

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



## JAVA版本 日志收集到腾讯云


### pom 导入
```
	<dependency>
	  <groupId>com.yondor.log</groupId>
	  <artifactId>applog</artifactId>
	  <version>0.0.3</version>
	</dependency>
```
### 增加 applog.properties 配置文件到classes文件夹:
```
project=AppLog
version=0.0.1
time=5000
appid=10
token=from_k8s
url=http\://apps-applog.default/logcat/server
```
### 说明：
```
project:工程名称
version:工程版本
time:发送日志频率
appid: 自定义appid ，大家不同就可以了。
token: 暂时没有用
url: 发布日志的url
```

### 用法：
```
public class XXXX{
     //声明全局变量
     final  AppLog applog=AppLog.getInstance();


    public void test(){

         applog.info("server_id=1","world",0L);// info tag,message,userID
         applog.warn("server_id=1","world",0L);// warn tag,message,userID
         applog.error("server_id=1","error message",0L);// error tag,message,userID
   }
}


```
