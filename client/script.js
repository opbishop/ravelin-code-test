window.addEventListener("load", initialiseVars);
window.addEventListener("resize", resize);
document.getElementById("submit").addEventListener("click", submitButton);
document.getElementById("inputEmail").addEventListener("keyup", handleTyping);
document.getElementById("inputEmail").addEventListener("paste", paste);
document.getElementById("inputCardNumber").addEventListener("keyup", handleTyping);
document.getElementById("inputCardNumber").addEventListener("paste", paste);
document.getElementById("inputCVV").addEventListener("keyup", handleTyping);
document.getElementById("inputCVV").addEventListener("paste", paste);

var firstCharTyped = 0;
var width;
var height;
const sessionId = create_UUID()
const serverAddr = "http://127.0.0.1:10000/"

function create_UUID(){
    var dt = new Date().getTime();
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (dt + Math.random()*16)%16 | 0;
        dt = Math.floor(dt/16);
        return (c=='x' ? r :(r&0x3|0x8)).toString(16);
    });
    return uuid;
}

function initialiseVars() {
    width = window.innerWidth;
    height = window.innerHeight;
}

function handleTyping(){
    firstCharTyped = new Date().getSeconds();
    console.log("typing started at " + firstCharTyped)

    document.querySelectorAll('.form-control').forEach(item => {
        item.removeEventListener('keyup', handleTyping);
    })
}
function sendRequest(url, content){
    console.log(content)
    let request = new XMLHttpRequest();

    // send async request
    request.open("POST", url);
    request.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    request.send(JSON.stringify(content))
}

function paste(){
    const fieldName = this.id;
    let content =
        {
            "eventType": "copypaste",
            "websiteUrl": window.location["href"],
            "sessionId": sessionId,
            "formId": fieldName,
            "pasted": true
        }

    const url = serverAddr + content["eventType"]
    sendRequest(url, content)

}

function resize(){
    let content =
        {
            "eventType": "resize",
            "websiteUrl": window.location["href"],
            "sessionId": sessionId,
            "oldWidth": width.toString(),
            "oldHeight": height.toString(),
            "newWidth":  window.innerWidth.toString(),
            "newHeight": window.innerHeight.toString()
        }

    initialiseVars()

    const url = serverAddr + content["eventType"]
    sendRequest(url, content)
}

function submitButton(){
    console.log(firstCharTyped)
    console.log(new Date().getSeconds())
    console.log(new Date().getSeconds() - firstCharTyped)
    let time = new Date().getSeconds() - firstCharTyped;
    let content =
        {
            "eventType": "timer",
            "websiteUrl": window.location["href"],
            "sessionId": sessionId,
            "time": time
        }
    console.log(content)

    const url = serverAddr + content["eventType"]
    sendRequest(url, content)
}