var pHappy = document.getElementById("progress-happy").style;
var pAngry = document.getElementById("progress-angry").style;
var pNeutral = document.getElementById("progress-neutral").style;
var pSad = document.getElementById("progress-sad").style;
var pFearful = document.getElementById("progress-fearful").style;
var pDisgusted = document.getElementById("progress-disgusted").style;
var pSUrprised = document.getElementById("progress-surprised").style;

var i = 0;
var increase = true;

function pUpdate(){
    pHappy.width = i+"%";
    pSad.width = i+"%";
    pAngry.width = i+"%";
    pNeutral.width = i+"%";
    if(increase){
        i++;
        if(i == 100){
            increase = false;
        }
    } else {
        i--;
        if(i == 0){
            increase = true;
        }
    }
}

let timerId = setInterval(pUpdate, 100);
