
import QtQuick 2.0
//import QtMultimedia 5.0
import QtMultimedia 5.0
import QtQuick.Controls 1.2
import QtQuick.Controls.Styles 1.2
import QtQuick.Layouts 1.1
import QtQuick.Dialogs 1.2
import QtQuick.Window 2.0
import '../js/util.js' as Util
ApplicationWindow{
    id:kPlayer
    width: 1024
    height: 650
    visible: true
    title: getVedioName(fd.fileUrl.toString())

    onClosing: {
        kPlayer.destroy()
    }
    //获取影音名称
    function getVedioName(str)
    {
        var url=fd.fileUrl.toString()
        var strList= []
        strList=url.split("/");
        var name=strList[strList.length-1];
        return name;
    }
    Column{
        Rectangle{
            id:screen
            color:"black"
            width:kPlayer.width
            height: kPlayer.height-50
            Image{
                id:img
                source: "./Images/KPlayer.png"
                anchors.fill: parent
            }
            MediaPlayer{
                id:player
                source: fd.fileUrl
                autoPlay: true
                volume: voice.value
                onError: {
                    console.log('player error', error, errorString)
                }
            }
            VideoOutput {
                id: vo
                anchors.fill: parent
                source: player
                transform: Rotation {
                    id: videoRotation
                    origin.x: vo.width/2
                    origin.y: vo.height / 2
                    angle: 0
                }
            }

        }
        Rectangle{
            id:control
            color:"#80202020"
            border.color: "gray"
            border.width: 1
            width:kPlayer.width
            height: 20
            Row{
                spacing: 10
                anchors.verticalCenter: parent.verticalCenter
                anchors.leftMargin: 10
                anchors.left: parent.left
                //调节播放速度
                Slider{
                    id:playPos
                    width: kPlayer.width*0.75
                    height: 10
                    maximumValue: player.duration
                    minimumValue: 0
                    value:player.position
                    anchors.verticalCenter: parent.verticalCenter
                    stepSize:1000
                    style: SliderStyle {
                        groove: Rectangle {
                            width: kPlayer.width*0.8
                            height: 8
                            color: "gray"
                            radius: 2
                            Rectangle {
                                anchors.left: parent.left
                                anchors.top: parent.top
                                anchors.bottom: parent.bottom
                                width: player.duration>0?parent.width*player.position/player.duration:0
                                color: "blue"
                            }
                        }
                        handle: Rectangle {
                            anchors.centerIn: parent
                            color: control.pressed ? "white" : "darkgray"
                            border.color: "gray"
                            border.width: 2
                            implicitWidth: 15
                            implicitHeight: 15
                            radius:7.5
                            Rectangle{
                                width: parent.width-8
                                height: width
                                radius: width/2
                                color: "blue"
                                anchors.centerIn: parent
                            }
                        }
                    }
                    //点击鼠标设置播放位置
                    MouseArea {
                        property int pos
                        anchors.fill: parent
                        onClicked: {
                            if (player.seekable)
                                pos = player.duration * mouse.x/parent.width
                            player.seek(pos)
                            playPos.value=pos;
                        }
                    }
                }
                Image{
                    width: 15
                    height: 15
                    source: "./Images/voice.png"
                    anchors.verticalCenter: parent.verticalCenter
                }
                //调节音量
                Slider{
                    id:voice
                    width: kPlayer.width*0.2
                    height: 10
                    value:player.volume
                    stepSize: 0.1
                    maximumValue: 1
                    minimumValue: 0
                    anchors.verticalCenter: parent.verticalCenter
                    style: SliderStyle {
                        groove: Rectangle {
                            implicitWidth: kPlayer.width*0.2
                            implicitHeight: 8
                            color: "gray"
                            radius: 2
                            Rectangle {
                                anchors.left: parent.left
                                anchors.top: parent.top
                                anchors.bottom: parent.bottom
                                width: player.volume>0?parent.width*player.volume:0
                                color: "blue"
                            }
                        }
                        handle: Rectangle {
                            anchors.centerIn: parent
                            color: control.pressed ? "white" : "darkgray"
                            border.color: "gray"
                            border.width: 2
                            implicitWidth: 15
                            implicitHeight: 15
                            radius:7.5
                            Rectangle{
                                width: parent.width-8
                                height: width
                                radius: width/2
                                color: "blue"
                                anchors.centerIn: parent
                            }
                        }
                    }
                }
            }
        }
        //控制区域
        Rectangle{
            id:bottomArea
            color:"#80202020"
            border.color: "gray"
            border.width: 1
            width: kPlayer.width
            height: 30
            Row{
                anchors.verticalCenter: parent.verticalCenter
                anchors.left: parent.left
                anchors.leftMargin: 10
                spacing: 10
                Button {
                   text: '旋转'
                   onClicked: {
                       videoRotation.angle  = (videoRotation.angle + 90) % 360
                   }
                }
                ButtonImage {
                    width: 30
                    height: 30
                    property int status: 1  //默认播放
                    iconImage: "./Images/pause.png"
                    onClicked: {
                        if(status===1)
                        {
                            player.pause();
                            tooltip="开始";
                            console.log("start")
                            status=0;
                            iconImage="./Images/play.png"
                        }
                        else{
                            player.play() ;
                            tooltip="暂停";
                            console.log("pause")
                            status=1;
                            iconImage="./Images/pause.png"
                        }
                    }
                }
                ButtonImage {
                    width: 30
                    height: 30
                    onClicked: player.stop()
                    tooltip: "停止"
                    iconImage: "./Images/stop.png"
                }
                //快进快退10s
                ButtonImage {
                    width: 30
                    height: 30
                    onClicked: player.seek(player.position+10000)
                    tooltip: "快退"
                    iconImage: "./Images/back.png"
                }
                ButtonImage {
                    width: 30
                    height: 30
                    onClicked: player.seek(player.position-10000)
                    tooltip: "快进"
                    iconImage: "./Images/pass.png"
                }
                ButtonImage{
                    width: 30
                    height: 30
                    tooltip: "打开文件"
                    onClicked: fd.open()
                    iconImage: "./Images/add.png"
                    FileDialog{
                        id:fd
                        nameFilters: ["Vedio Files(*.avi *.mp4 *rmvb *.rm)"]  //格式过滤
                        selectMultiple: false
                    }
                }
                Text{
                    anchors.verticalCenter: parent.verticalCenter
                    text:parent.currentTime(player.position)+"/"+parent.currentTime(player.duration)
                    color: "white"
                }
                //时间格式化
                function currentTime(time)
                {
                    var sec= Math.floor(time/1000);
                    var hours=Math.floor(sec/3600);
                    var minutes=Math.floor((sec-hours*3600)/60);
                    var seconds=sec-hours*3600-minutes*60;
                    var hh,mm,ss;
                    if(hours.toString().length<2)
                        hh="0"+hours.toString();
                    else
                        hh=hours.toString();
                    if(minutes.toString().length<2)
                        mm="0"+minutes.toString();
                    else
                        mm=minutes.toString();
                    if(seconds.toString().length<2)
                        ss="0"+seconds.toString();
                    else
                        ss=seconds.toString();
                    return hh+":"+mm+":"+ss
                }
            }
        }
    }
    function playVideo(meta, useVip) {
        player.stop()
        title = meta.server_filename
        Util.callGoAsync('pan.link', {fid: meta.fs_id, useVip: useVip})
        .then(function(link) {
            var agentLink = Util.videoAgentLink(meta, useVip)
            console.log('play link', agentLink, link)
            player.source = agentLink
        })
    }

}
