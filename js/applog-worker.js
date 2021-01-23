'use strict';

var init = false;
var head = {};
var dataList = [];
const MAX_SIZE = 200;
var CHECK_TIME = 10000;

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

function checkMessage() {

    postMessage(dataList.length);
    autoSubmit();
    setTimeout("checkMessage()", CHECK_TIME);
}

function autoSubmit() {

    if (dataList.length == 0) {
        // console.log("@autoSubmit dataList.length=",dataList.length)
        return;
    }
    if (!init) {
        console.log("@autoSubmit error: applog-worker do not init yet!")
        return;
    }

    var datas = dataList;
    if (dataList.length > MAX_SIZE) {
        datas = dataList.slice(0, MAX_SIZE);
    }
    dataList = [];
    var data = {
        "appid": head.appid,
        "version": head.version,
        "uid": head.uid,
        "device": head.device,
        "platform": head.platform,
        "key": head.uuid,
        "data": datas
    };
    var url = head.url + '/logcat/app';

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json; charset=UTF-8'
        },
        body: JSON.stringify(data)
    }).then(checkStatus).then(parseJSON).then(function(ret) {
        // console.log("applog submit ret ",ret);

    }).catch(function(err) {
        console.error(err);
        dataList = dataList.concat(datas);
    });

}
self.onmessage = function(e) {
    var d = e.data;
    if (d.type === "init") {
        head = d;
        init = true;
        return;
    } else if (e.type === "uid") {
        head.userID = d.userID;
        return;
    }
    dataList.push(d.data);
};
postMessage("init");
checkMessage();