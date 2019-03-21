#!/bin/bash
set -ev

#check env
df -h

ls $HOME/*
du -sh $HOME/*

#needed for headless qt installation
export QT_QPA_PLATFORM=minimal

#additional deps for multimedia
sudo apt-get -y -qq install libpulse-dev && sudo apt-get -qq clean

if [ "$QT_PKG_CONFIG" == "true" ]
then
  #download and install qt
  sudo add-apt-repository -y ppa:forkotov02/opt-qt-5.10.0-trusty
  sudo apt-get -qq update
  sudo apt-get -y -qq install qt510base qt510connectivity qt510declarative qt510doc qt510location qt510multimedia qt510quickcontrols qt510quickcontrols2 qt510script qt510serialport qt510svg qt510tools qt510translations qt510webchannel qt510webengine qt510websockets qt510x11extras qt510xmlpatterns && sudo apt-get -qq clean
else
  #download and install qt
  QT=qt-unified-linux-x64-online.run
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$QT https://download.qt.io/official_releases/online_installers/$QT
  chmod +x /tmp/$QT
  /tmp/$QT -v --script $GOPATH/src/github.com/peterq/pan-light/qt/internal/ci/iscript.qs LINUX=true
  rm -f /tmp/$QT
fi

if [ "$QT_MXE" == "true" ]
then
  #download and install qt (and wine) for cross compilation
  sudo apt-get -y -qq install wine && sudo apt-get -qq clean
  echo "deb http://pkg.mxe.cc/repos/apt/debian wheezy main" | sudo tee --append /etc/apt/sources.list.d/mxeapt.list > /dev/null
  sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys D43A795B73B16ABE9643FE1AFD8FFF16DB45C6AB
  sudo apt-get -qq update

  if [ "$QT_MXE_STATIC" == "true" ]
  then
    if [ "$QT_MXE_ARCH" == "386" ]
    then
      sudo apt-get -y -qq install mxe-i686-w64-mingw32.static-qt3d mxe-i686-w64-mingw32.static-qtactiveqt mxe-i686-w64-mingw32.static-qtbase mxe-i686-w64-mingw32.static-qtcanvas3d mxe-i686-w64-mingw32.static-qtcharts mxe-i686-w64-mingw32.static-qtconnectivity mxe-i686-w64-mingw32.static-qtdatavis3d mxe-i686-w64-mingw32.static-qtdeclarative mxe-i686-w64-mingw32.static-qtgamepad mxe-i686-w64-mingw32.static-qtgraphicaleffects mxe-i686-w64-mingw32.static-qtimageformats mxe-i686-w64-mingw32.static-qtlocation mxe-i686-w64-mingw32.static-qtmultimedia mxe-i686-w64-mingw32.static-qtofficeopenxml mxe-i686-w64-mingw32.static-qtpurchasing mxe-i686-w64-mingw32.static-qtquickcontrols mxe-i686-w64-mingw32.static-qtquickcontrols2 mxe-i686-w64-mingw32.static-qtscript mxe-i686-w64-mingw32.static-qtscxml mxe-i686-w64-mingw32.static-qtsensors mxe-i686-w64-mingw32.static-qtserialbus mxe-i686-w64-mingw32.static-qtserialport mxe-i686-w64-mingw32.static-qtservice mxe-i686-w64-mingw32.static-qtsvg mxe-i686-w64-mingw32.static-qtsystems mxe-i686-w64-mingw32.static-qttools mxe-i686-w64-mingw32.static-qttranslations mxe-i686-w64-mingw32.static-qtvirtualkeyboard mxe-i686-w64-mingw32.static-qtwebchannel mxe-i686-w64-mingw32.static-qtwebsockets mxe-i686-w64-mingw32.static-qtwinextras mxe-i686-w64-mingw32.static-qtxlsxwriter mxe-i686-w64-mingw32.static-qtxmlpatterns
    else
      sudo apt-get -y -qq install mxe-x86-64-w64-mingw32.static-qt3d mxe-x86-64-w64-mingw32.static-qtactiveqt mxe-x86-64-w64-mingw32.static-qtbase mxe-x86-64-w64-mingw32.static-qtcanvas3d mxe-x86-64-w64-mingw32.static-qtcharts mxe-x86-64-w64-mingw32.static-qtconnectivity mxe-x86-64-w64-mingw32.static-qtdatavis3d mxe-x86-64-w64-mingw32.static-qtdeclarative mxe-x86-64-w64-mingw32.static-qtgamepad mxe-x86-64-w64-mingw32.static-qtgraphicaleffects mxe-x86-64-w64-mingw32.static-qtimageformats mxe-x86-64-w64-mingw32.static-qtlocation mxe-x86-64-w64-mingw32.static-qtmultimedia mxe-x86-64-w64-mingw32.static-qtofficeopenxml mxe-x86-64-w64-mingw32.static-qtpurchasing mxe-x86-64-w64-mingw32.static-qtquickcontrols mxe-x86-64-w64-mingw32.static-qtquickcontrols2 mxe-x86-64-w64-mingw32.static-qtscript mxe-x86-64-w64-mingw32.static-qtscxml mxe-x86-64-w64-mingw32.static-qtsensors mxe-x86-64-w64-mingw32.static-qtserialbus mxe-x86-64-w64-mingw32.static-qtserialport mxe-x86-64-w64-mingw32.static-qtservice mxe-x86-64-w64-mingw32.static-qtsvg mxe-x86-64-w64-mingw32.static-qtsystems mxe-x86-64-w64-mingw32.static-qttools mxe-x86-64-w64-mingw32.static-qttranslations mxe-x86-64-w64-mingw32.static-qtvirtualkeyboard mxe-x86-64-w64-mingw32.static-qtwebchannel mxe-x86-64-w64-mingw32.static-qtwebsockets mxe-x86-64-w64-mingw32.static-qtwinextras mxe-x86-64-w64-mingw32.static-qtxlsxwriter mxe-x86-64-w64-mingw32.static-qtxmlpatterns
    fi
  else
    if [ "$QT_MXE_ARCH" == "386" ]
    then
      sudo apt-get -y -qq install mxe-i686-w64-mingw32.shared-qt3d mxe-i686-w64-mingw32.shared-qtactiveqt mxe-i686-w64-mingw32.shared-qtbase mxe-i686-w64-mingw32.shared-qtcanvas3d mxe-i686-w64-mingw32.shared-qtcharts mxe-i686-w64-mingw32.shared-qtconnectivity mxe-i686-w64-mingw32.shared-qtdatavis3d mxe-i686-w64-mingw32.shared-qtdeclarative mxe-i686-w64-mingw32.shared-qtgamepad mxe-i686-w64-mingw32.shared-qtgraphicaleffects mxe-i686-w64-mingw32.shared-qtimageformats mxe-i686-w64-mingw32.shared-qtlocation mxe-i686-w64-mingw32.shared-qtmultimedia mxe-i686-w64-mingw32.shared-qtofficeopenxml mxe-i686-w64-mingw32.shared-qtpurchasing mxe-i686-w64-mingw32.shared-qtquickcontrols mxe-i686-w64-mingw32.shared-qtquickcontrols2 mxe-i686-w64-mingw32.shared-qtscript mxe-i686-w64-mingw32.shared-qtscxml mxe-i686-w64-mingw32.shared-qtsensors mxe-i686-w64-mingw32.shared-qtserialbus mxe-i686-w64-mingw32.shared-qtserialport mxe-i686-w64-mingw32.shared-qtservice mxe-i686-w64-mingw32.shared-qtsvg mxe-i686-w64-mingw32.shared-qtsystems mxe-i686-w64-mingw32.shared-qttools mxe-i686-w64-mingw32.shared-qttranslations mxe-i686-w64-mingw32.shared-qtvirtualkeyboard mxe-i686-w64-mingw32.shared-qtwebchannel mxe-i686-w64-mingw32.shared-qtwebkit mxe-i686-w64-mingw32.shared-qtwebsockets mxe-i686-w64-mingw32.shared-qtwinextras mxe-i686-w64-mingw32.shared-qtxlsxwriter mxe-i686-w64-mingw32.shared-qtxmlpatterns
    else
      sudo apt-get -y -qq install mxe-x86-64-w64-mingw32.shared-qt3d mxe-x86-64-w64-mingw32.shared-qtactiveqt mxe-x86-64-w64-mingw32.shared-qtbase mxe-x86-64-w64-mingw32.shared-qtcanvas3d mxe-x86-64-w64-mingw32.shared-qtcharts mxe-x86-64-w64-mingw32.shared-qtconnectivity mxe-x86-64-w64-mingw32.shared-qtdatavis3d mxe-x86-64-w64-mingw32.shared-qtdeclarative mxe-x86-64-w64-mingw32.shared-qtgamepad mxe-x86-64-w64-mingw32.shared-qtgraphicaleffects mxe-x86-64-w64-mingw32.shared-qtimageformats mxe-x86-64-w64-mingw32.shared-qtlocation mxe-x86-64-w64-mingw32.shared-qtmultimedia mxe-x86-64-w64-mingw32.shared-qtofficeopenxml mxe-x86-64-w64-mingw32.shared-qtpurchasing mxe-x86-64-w64-mingw32.shared-qtquickcontrols mxe-x86-64-w64-mingw32.shared-qtquickcontrols2 mxe-x86-64-w64-mingw32.shared-qtscript mxe-x86-64-w64-mingw32.shared-qtscxml mxe-x86-64-w64-mingw32.shared-qtsensors mxe-x86-64-w64-mingw32.shared-qtserialbus mxe-x86-64-w64-mingw32.shared-qtserialport mxe-x86-64-w64-mingw32.shared-qtservice mxe-x86-64-w64-mingw32.shared-qtsvg mxe-x86-64-w64-mingw32.shared-qtsystems mxe-x86-64-w64-mingw32.shared-qttools mxe-x86-64-w64-mingw32.shared-qttranslations mxe-x86-64-w64-mingw32.shared-qtvirtualkeyboard mxe-x86-64-w64-mingw32.shared-qtwebchannel mxe-x86-64-w64-mingw32.shared-qtwebkit mxe-x86-64-w64-mingw32.shared-qtwebsockets mxe-x86-64-w64-mingw32.shared-qtwinextras mxe-x86-64-w64-mingw32.shared-qtxlsxwriter mxe-x86-64-w64-mingw32.shared-qtxmlpatterns
    fi
  fi

  sudo apt-get -qq clean
fi

if [ "$ANDROID" == "true" ]; then
  #install openjdk8
  echo "export JDK_DIR=/usr/lib/jvm/java-8-openjdk-amd64" >> $HOME/.profile
  sudo add-apt-repository -y ppa:openjdk-r/ppa
  sudo apt-get -qq update
  sudo apt-get -y -qq install openjdk-8-jdk && sudo apt-get -qq clean

  #download and install android sdk
  SDK=sdk-tools-linux-3859397.zip
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$SDK https://dl.google.com/android/repository/$SDK
  unzip -qq /tmp/$SDK -d $HOME/android-sdk-linux/
  rm -f /tmp/$SDK

  #install deps for android sdk
  $HOME/android-sdk-linux/tools/bin/sdkmanager --list --verbose
  echo "y" | $HOME/android-sdk-linux/tools/bin/sdkmanager "platform-tools" "build-tools;26.0.0" "platforms;android-25"
  $HOME/android-sdk-linux/tools/bin/sdkmanager --update

  #download and install android ndk
  NDK=android-ndk-r18b-linux-x86_64.zip
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$NDK https://dl.google.com/android/repository/$NDK
  unzip -qq /tmp/$NDK -d $HOME
  rm -f /tmp/$NDK
fi

#prepare env
sudo chown $USER /usr/local/bin
sudo chown $USER $GOROOT/pkg | true

#check env
df -h

ls $HOME/*
du -sh $HOME/*

exit 0
