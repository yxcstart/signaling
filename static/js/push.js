'use strict';

var localVideo = document.getElementById("localVideo");
var pushBtn = document.getElementById("pushBtn");

pushBtn.addEventListener("click", startPush);

var uid = $("#uid").val();
var streamName = $("#streamName").val();
var audio = $("#audio").val();
var video = $("#video").val();
var offer = "";
var pc;
const config = {};
var localStream;
var lastConnectionState = "";

function startPush() {
    console.log("send push: /signaling/push");

    $.post("/signaling/push",
        { "uid": uid, "streamName": streamName, "audio": audio, "video": video },
        function (data, textStatus) {
            console.log("push response: " + JSON.stringify(data))
            if ("success" == textStatus && 0 == data.code) {
                $("#tips1").html("<font color='blue'>推流请求成功!</font>");
                console.log("remote offer: \r\n" + data.data.sdp);
                offer = data.data;
                pushStream();
            } else {
                $("#tips1").html("<font color='red'>推流请求失败!</font>");
            }
        },
        "json"
    );
}

function sendAnswer(answerSdp) {
    console.log("send answer: /signaling/sendanswer");

    $.post("/signaling/sendanswer",
        { "uid": uid, "streamName": streamName, "answer": answerSdp, "type": "push" },
        function (data, textStatus) {
            console.log("send answer response: " + JSON.stringify(data));
            if ("success" == textStatus && 0 == data.code) {
                $("#tips3").html("<font color='blue'>answer发送成功!</font>");
            } else {
                $("#tips3").html("<font color='red'>answer发送失败!</font>");
            }
        },
        "json"
    );
}


function pushStream() {
    pc = new RTCPeerConnection(config);

    pc.oniceconnectionstatechange = function (e) {
        var state = "";
        if (lastConnectionState != "") {
            state = lastConnectionState + "->" + pc.iceConnectionState;
        } else {
            state = pc.iceConnectionState;
        }

        $("#tips2").html("连接状态: " + state);
        lastConnectionState = pc.iceConnectionState;
    }

    pc.setRemoteDescription(offer).then(
        setRemoteDescriptionSuccess,
        setRemoteDescriptionError
    );
}

window.addEventListener("message", function (event) {
    if (event.origin != window.location.origin) {
        return;
    }

    if (event.data.type) {
        if (event.data.type == "SS_DIALOG_SUCCESS") {
            console.log("用户同意屏幕共享, streamId: " + event.data.streamId);
            startScreenStreamFrom(event.data.streamId);
        } else if (event.data.type == "SS_DIALOG_CANCEL") {
            console.log("用户取消屏幕共享");
        }
    }
});

function setRemoteDescriptionSuccess() {
    console.log("pc set remote description success");
    console.log("request screen share");
    window.postMessage({ type: "SS_UI_REQUEST", text: "push" }, "*");
}

function setRemoteDescriptionError(error) {
    console.log("pc set remote description error: " + error);
}

function startScreenStreamFrom(streamId) {
    var constrains = {
        audio: false,
        video: {
            mandatory: {
                chromeMediaSource: "desktop",
                chromeMediaSourceId: streamId,
                maxWidth: window.screen.width,
                maxHeight: window.screen.height
            }
        }
    };

    navigator.mediaDevices.getUserMedia(constrains).then(
        handleSuccess).catch(handleError);
}

function handleSuccess(stream) {
    navigator.mediaDevices.getUserMedia({ audio: true }).then(
        function (audioStream) {
            stream.addTrack(audioStream.getAudioTracks()[0]);
            localVideo.srcObject = stream;
            localStream = stream;
            pc.addStream(stream);
            pc.createAnswer().then(
                createSessionDescriptionSuccess,
                createSessionDescriptionError
            );
        }
    ).catch(handleError);
}

function handleError(error) {
    console.log("get user media error: " + error);
}

function createSessionDescriptionSuccess(answer) {
    console.log("answer sdp: \n" + answer.sdp);
    console.log("pc set local sdp");
    pc.setLocalDescription(answer).then(
        setLocalDescriptionSuccess,
        setLocalDescriptionError
    );
    sendAnswer(answer.sdp);
}

function setLocalDescriptionSuccess() {
    console.log("set local description success");
}

function setLocalDescriptionError(error) {
    console.log("pc set local description error: " + error);
}

function createSessionDescriptionError(error) {
    console.log("pc create answer error: " + error);
}