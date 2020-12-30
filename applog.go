package applog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/widuu/goini"
)

// 提交日志到腾讯云

//default URL
var defaultURL="http://apps-applog.default/logcat/server"

var hasInit=false
var lastWarn=""
var confData = ConfigData{}
var revChan chan string= make(chan string,1000)
var head1 string
var head2 string

//LevelFine LevelFine
const LevelFine=0
//LevelInfo LevelInfo
const LevelInfo=1
//LevelWarn LevelWarn
const LevelWarn=2
//LevelError LevelError
const LevelError=3

//ConfigData ConfigData
type ConfigData struct{
	Project	 	string
	Version		string
	Time 			int
	Appid			int
	Token			string
	URL 			string
	Enable		bool
	Level			int
}

// AppConfig AppConfig
type AppConfig struct{
	Code   int `json:"code"`
	Enable bool `json:"enable"`
	Level  int `json:"level"`
}

//confWarn confWarn
func confWarn(msg string){
	log.Println("==========Log Config Warn===========")
	if msg!=""{
		log.Println("Error:",msg)
	}else{
		log.Println(lastWarn)
	}
	lastWarn=msg
}

// init init
func init(){
	hasInit=false
	appIniFile:=""
	if fileExists("./data/config/app.ini"){
		appIniFile="./data/config/app.ini"
	}else	if fileExists("app.ini"){
		appIniFile="app.ini"
	}else	if fileExists("log.ini"){
		appIniFile="log.ini"
	}else{
		confWarn("Error:please config data/config/app.ini or app.ini or log.ini for log!")
		return
	}

	conf := goini.SetConfig(appIniFile)

	confProject:=conf.GetValue("log", "project")
	confVersion:=conf.GetValue("log", "version")
	confTime:=conf.GetValue("log", "time")
	confAppid:=conf.GetValue("log", "appid")
	confToken:=conf.GetValue("log", "token")
	confURL:=conf.GetValue("log", "url")

	if confProject=="no value"{
		confWarn(appIniFile+": Please config 'project' in [log]")
		return
	}
	if confVersion=="no value"{
		confWarn(appIniFile+": Please config 'version' in [log]")
		return
	}
	if confTime=="no value"{
		confWarn(appIniFile+": Please config 'time' in [log]")
		return
	}
	if confAppid=="no value"{
		confWarn(appIniFile+": Please config 'appid' in [log]")
		return
	}
	if confToken=="no value"{
		confWarn(appIniFile+": Please config 'token' in [log]")
		return
	}
	if confURL=="no value"{
		confURL=defaultURL
	}
	if confURL==""{
		confURL=defaultURL
	}

	var err error=nil
	confData.Appid,err=strconv.Atoi(confAppid)
	if err!=nil{
		confWarn(appIniFile+": 'appid' is wrong,err:"+err.Error())
		return
	}
	confData.Project=confProject
	confData.Time,err=strconv.Atoi(confTime)
	if err!=nil{
		confWarn(appIniFile+": 'time' is wrong,err:"+err.Error())
		return
	}
	if confData.Time<1000{
		confData.Time=1000
		confWarn(appIniFile+": 'time' is too small,please set time >=1000 ")
		return
	}
	
	confData.Token=confToken
	confData.URL=confURL
	confData.Version=confVersion
	confData.Enable=true
	confData.Level=LevelFine
 
	_uuid,_:=uuid.NewRandom()

	key:=_uuid.String()
	key=strings.ReplaceAll(key,"-","")
	

	head1="{\"appid\":"+confAppid+",\"version\":\""+filter(confVersion)+"\",\"key\":\""+filter(key)+"\",\"name\":\""+filter(confProject)+"\",\"token\":\""+filter(confToken)+"\",\"data\":";
	head2="}";
 

	go revLog()
  go checkEnable()

	hasInit=true
}
func revLog(){
	// revChan = make(chan string,1000)
	data := make([]string,1000)//固定数组
 	i:=0
	nextTime:=time.Now().Add(time.Duration(confData.Time)*time.Millisecond).Unix()
	for{
		//收集数据
		msg:=<-revChan
		if i<1000{
			data[i]=msg
		  i++
		}else{
			//ignore
			//log.Println("applog@revLog 忽略数据")
		}
		//发送数据
		cur:=time.Now().Unix()
		if cur>=nextTime && i>0 {
			//发送
			nextTime=time.Now().Add(time.Duration(confData.Time)*time.Millisecond).Unix()
			buf:=bytes.Buffer{}
			buf.WriteString(head1)
			buf.WriteString("[")
			for j:=0;j<i;j++{
				if j>0{
					buf.WriteString(",")
				}
				buf.WriteString(data[j])
				
			}
			buf.WriteString("]")
			buf.WriteString(head2)

			if post(buf.String()){
				//清空旧的数据
				for j:=0;j<i;j++{
					data[j]=""
				}
				i=0
			}
		}
	}
}
//LevelName LevelName
func LevelName(level int) string{
	if level==0{
		return "fine"
	}
	if level==1{
		return "info"
	}
	if level==2{
		return "warn"
	}
	if level==3{
		return "error"
	}
	return "unknow"
}
// checkEnable checkEnable
func checkEnable(){
	client := &http.Client{}
	part:=fmt.Sprintf("logcat/config/%d/1",confData.Appid)
	checkURL:=strings.Replace(confData.URL,"logcat/server",part,-1)

	req, err := http.NewRequest("GET",checkURL,nil)
	if err!=nil{
		log.Println("applog@checkEnable err:",err.Error())
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if resp.StatusCode<200 || resp.StatusCode>=300{
		log.Println("applog@checkEnable status code:",resp.StatusCode )
		return
	}
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		log.Println("applog@checkEnable err:",errRead.Error())
		return  
	}
	
	appconfig :=AppConfig{}
	errJSON:=json.Unmarshal(body,&appconfig)
	if errJSON != nil {
		log.Println("applog@checkEnable errJSON:",errJSON.Error())
		return  
	}
	if appconfig.Code>=200 && appconfig.Code<300{
		confData.Enable=appconfig.Enable
		confData.Level=appconfig.Level
	}
	if confData.Enable==false{

	}

	log.Println("applog@checkEnable Enable:",confData.Enable)
	log.Println("applog@checkEnable Level:",LevelName(confData.Level))

}
//post post
func post(data string) bool{

	client := &http.Client{}
	log.Println("@post ",confData.URL,data)
  req, err := http.NewRequest("POST", confData.URL, strings.NewReader(data))
  if err != nil {
			// handle error
			log.Println("applog@post ",err.Error())
			return false
  }
	req.Header.Set("Content-Type", "application/json")
	
  // req.Header.Set("Cookie", "name=anny")
 
	resp, errRequest := client.Do(req)
	if errRequest!=nil{
		log.Println("applog@post err:",errRequest.Error())
		return false
	}
  defer resp.Body.Close()
	if resp.StatusCode>=200 && resp.StatusCode<300 {
		//ok
		return true
	}
  body, errRead := ioutil.ReadAll(resp.Body)
  if errRead != nil {
			log.Println("applog@post err:",errRead.Error())
			return false
	}
	content:=string(body)
	fmt.Println("applog@post message:",content)
	return false
}
 
// 判断资源包是否存在
func fileExists(filename string) (bool) {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	return false
}

// Fine Fine
func Fine(tag string,message string,uid int64) {
	_log(LevelFine,tag,message,uid)
}

// Info Info
func Info(tag string,message string,uid int64) {
	_log(LevelInfo,tag,message,uid)
}

// Warn Warn
func Warn(tag string,message string,uid int64) {
	_log(LevelWarn,tag,message,uid)
}

// Error Error
func Error(tag string,message string,uid int64) {
	_log(LevelError,tag,message,uid)
}

// _log _log
func _log(level int,tag string,message string,uid int64) {
	 if !hasInit{
		 confWarn("")
		 if level>=LevelInfo{
			 fmt.Println(tag,message,uid)
		 }
		 return
	 }

	 if confData.Enable==false{
		 //没有生效
		 return
	 }
	 timestamp:=time.Now().Unix()
 
	 rownum:=1
	 file:="main.go"

	 d:= fmt.Sprintf("{\"t\":%d,\"l\":%d,\"g\":\"%s\",\"c\":\"%s\",\"u\":%d,\"s\":\"%s\",\"r\":%d}",timestamp,level,tag,filter(message),uid,filter(file),rownum)
	 revChan<-d
}
func filter(str string) string{
	if len(str)==0 {
		return str
	}
	str=strings.Join(strings.Split(str,"\""),"\\\"")
	str=strings.Join(strings.Split(str,"\r"),"\\r")
	str=strings.Join(strings.Split(str,"\n"),"\\n")
	str=strings.Join(strings.Split(str,"\t"),"\\t")
	str=strings.Join(strings.Split(str,"\f"),"\\f")
	str=strings.Join(strings.Split(str,"\b"),"\\b")
	return str
}
