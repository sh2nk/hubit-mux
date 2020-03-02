var webSocket = new WebSocket("ws://"+window.location.host+"/ws");

var pHappy = document.getElementById("progress-happy").style;
var pAngry = document.getElementById("progress-angry").style;
var pNeutral = document.getElementById("progress-neutral").style;
var pSad = document.getElementById("progress-sad").style;
var pFearful = document.getElementById("progress-fearful").style;
var pDisgusted = document.getElementById("progress-disgusted").style;
var pSurprised = document.getElementById("progress-surprised").style;

var i = 0;
var increase = true;

function pUpdate(){
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            console.log(this.responseText)
            var faces = JSON.parse(this.responseText);

            pSurprised.width = faces.emotions.surprised * 100 +"%";
            pAngry.width = faces.emotions.angry * 100 +"%";
            pHappy.width = faces.emotions.happy * 100 +"%";
            pNeutral.width = faces.emotions.neutral * 100 +"%";
            pFearful.width = faces.emotions.fearful * 100 +"%";
            pSad.width = faces.emotions.sad * 100 +"%"
            pDisgusted.width = faces.emotions.disgusted * 100 +"%";
        }
    };
    xmlhttp.open("GET", "/getface", true);
    xmlhttp.send();
}

let timerId = setInterval(pUpdate, 100);