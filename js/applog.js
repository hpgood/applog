'use strict';

// applog version:1.0.0 2021-01-20/Gance
//
var applog = (function(window) {
  const FINE=0;
  const INFO=1;
  const WARN=2;
  const ERROR=3;
  const uuidKey='applog_uuid';
  var appID=0;
  var serverURL="https://api.yondor.cn";
  var userID=0;
  var support=false;
  var enable=true;
  var uuid="";
  var appName="";
  var device="";
  var platform="";
  var worker=null;
  var workerJS="js/applog-worker.js";
  var curLevel=INFO;
  var _console_log=console.log;
  var _console_debug=console.debug;
  var _console_warn=console.warn;
  var _console_error=console.error;
  
 
  function parseJSON(response) {
    return response.json();
  }
  function checkStatus(response) {
    if (response.status >= 200 && response.status < 300) {
      return response 
    } else {
      var error = new Error(response.statusText);
      error.response = response
      throw error
    }
  }
  function guid2() {
    function S4() {
      return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1);
    }
    return (S4() + S4() + S4() +  S4() + S4() + + S4() + S4() + S4());
  }
  function getDevice(){
    var su = navigator.userAgent.toLowerCase();
    var  mb = ['ipad', 'iphone os', 'midp', 'rv:1.2.3.4', 'ucweb', 'android', 'windows ce', 'windows mobile',"ios","safari","mqqbrowser","chrome","mac","win"];
    for (var i in mb) {
      if (su.indexOf(mb[i]) > -1) {
        return mb[i];
      }
    }
    return "unknow";
  }
  function join2(arr){
    var str="";
    for(var i=1;i<arr.length;i++){
      if(i>1){
        str+=" ";
      }
      str+=arr[i];
    }
  }
  function getOS() {
    var sUserAgent = navigator.userAgent;
    var isWin = (navigator.platform == "Win32") || (navigator.platform == "Windows");
    var isMac = (navigator.platform == "Mac68K") || (navigator.platform == "MacPPC") || (navigator.platform == "Macintosh") || (navigator.platform == "MacIntel");
    if (isMac) return "Mac";
    var isUnix = (navigator.platform == "X11") && !isWin && !isMac;
    if (isUnix) return "Unix";
    var isLinux = (String(navigator.platform).indexOf("Linux") > -1);
    if (isLinux) return "Linux";
    if (isWin) {
      var isWin2K = sUserAgent.indexOf("Windows NT 5.0") > -1 || sUserAgent.indexOf("Windows 2000") > -1;
      if (isWin2K) return "Win2000";
      var isWinXP = sUserAgent.indexOf("Windows NT 5.1") > -1 || sUserAgent.indexOf("Windows XP") > -1;
      if (isWinXP) return "WinXP";
      var isWin2003 = sUserAgent.indexOf("Windows NT 5.2") > -1 || sUserAgent.indexOf("Windows 2003") > -1;
      if (isWin2003) return "Win2003";
      var isWinVista= sUserAgent.indexOf("Windows NT 6.0") > -1 || sUserAgent.indexOf("Windows Vista") > -1;
      if (isWinVista) return "WinVista";
      var isWin7 = sUserAgent.indexOf("Windows NT 6.1") > -1 || sUserAgent.indexOf("Windows 7") > -1;
      if (isWin7) return "Win7";
      var isWin10 = sUserAgent.indexOf("Windows NT 10") > -1 || sUserAgent.indexOf("Windows 10") > -1;
      if (isWin10) return "Win10";
    }
    return "other";
  }
  function startWorker(){
    if(worker!==null){
      return;
    }
    var data= {
      "type":"init",
      "appid":appID,  
      "version":appName, 
      "uid":1,
      "device":device,
      "platform":platform,
      "key":uuid,
      "url":serverURL
    };

    worker=new Worker(workerJS);
    worker.onerror=function(e){
      console.log("@onerror",e);
    }
    worker.onmessage=function(e){
      // console.log("@onmessage",e.data);
      if(e.data=="init"){
        // str=JSON.stringify(data);
        worker.postMessage(data);
      }
    }
    // _console_log("@startWorker worker=",worker);
  }

  return {
    init:function(id,name,workerPath){
      try{
        appID=id;
        userID=0;
        appName=name;
        device=getDevice();
        platform=getOS();
        if(typeof localStorage !=='undefined'){
          var last_uuid=localStorage.getItem(uuidKey);
          if(last_uuid && last_uuid!==""){
            uuid=last_uuid;
          }else{
            uuid=guid2();
            localStorage.setItem(uuidKey,uuid)
          }
        }else{
          uuid=guid2();
        }
        
        workerJS=workerPath
        if (typeof Worker ==="undefined" || typeof fetch ==="undefined"|| typeof JSON ==="undefined"){
          support=false;
          _console_error("@init your browser do not support fetcth/Worker api!");
        }else{
          support=true;
        }
        startWorker();
        this.fetchConfig();
      }catch(err){
        _console_error(err);
      }
    },
    setUserID:function(uid){
      try{
        if(typeof uid==='string'){
          userID=parseInt(uid);
        }else{
          userID=uid;
        }
        if(worker!=null){
          worker.postMessage({"type":"uid","userID":userID});
        }
      }catch(err){
        _console_error(err);
      }
    },
    setServerURL:function(url){
      serverURL=url;
    },
    setRedirect:function(t){
      try{
        var that=this;
        if(support && t){
          console.debug=function(){
            if(arguments.length==1){
              that.fine("debug",""+arguments[0]);
            }else if(arguments.length>=2){
              that.fine(""+arguments[0],join2(arguments));
            }
            
          };
          console.log=function(tag,event){
            if(arguments.length==1){
              that.fine("info",""+arguments[0]);
            }else if(arguments.length>=2){
              that.fine(""+arguments[0],join2(arguments));
            }
          };
          console.warn=function(tag,event){
            if(arguments.length==1){
              that.fine("warn",""+arguments[0]);
            }else if(arguments.length>=2){
              that.fine(""+arguments[0],join2(arguments));
            }
          };
          console.error=function(tag,event){
            if(arguments.length==1){
              that.fine("error",""+arguments[0]);
            }else if(arguments.length>=2){
              that.fine(""+arguments[0],join2(arguments));
            }
          };
        }else if(!t){
          console.debug=_console_debug;
          console.log=_console_log;
          console.warn=_console_warn;
          console.error=_console_error;
        }
      }catch(err){
        _console_error(err);
      }
    },
    fetchConfig: function(){
      if(!support){
        return;
      }
      fetch(serverURL+'/logcat/config/'+appID+"/"+userID).then(checkStatus).then(parseJSON).then(function(data){
        _console_log("applog config:",data);
        if(data.code>=200 && data.code<300){
          enable=data.enable;
          curLevel=data.level;
        }
      }).catch(function(err){
        console.error(err);
      });
    },
    fine:function(tag,message){
      this._log(FINE,tag,message);
    },
    info:function(tag,message){
      this._log(INFO,tag,message);
    },
    log:function(tag,message){
      this._log(INFO,tag,message);
    },
    warn:function(tag,message){
      this._log(WARN,tag,message);
    },
    error:function(tag,message){
      this._log(ERROR,tag,message);
    },
    _log:function(level,tag,message){
      try{
        if(level<curLevel){
          return;
        }
        var _tag=""+tag;
        var _msg=""+message;
        
        if(_tag.length>100){
          _tag=_msg.substring(0,100);
        }
        if(_msg.length>200){
          _msg=_msg.substring(0,200);
        }
        if(!enable || !support){
          switch(level){
            case FINE:
              console.debug(tag,message);
              break;
            case INFO:
              console.info(tag,message);
              break;
            case WARN:
              console.warn(tag,message);
              break;
            case ERROR:
              console.error(tag,message);
              break;
            default:
              console.info(tag,message);
          }
          return;
        }
        var timestamp=(new Date()).getTime();
        worker.postMessage({"type":"message","data":{"t":timestamp,"l":level,"g":tag,"c":message}});
      }catch(err){
        _console_error(err);
      }
    }
  }
})(window);
window.applog=applog;