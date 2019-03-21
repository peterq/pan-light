#!/bin/bash
set -ev

#check env
df -h
diskutil list

ls $HOME/*
du -sh $HOME/*

if [ "$QT_HOMEBREW" == "true" ]
then
  #download and install qt with brew
  brew update
  brew install qt5
  brew outdated qt5 || brew upgrade qt5
  ln -s /usr/local/Cellar/qt/5.11.2 $HOME/Desktop/Qt5.11.2
else
  #download and install qt
  QT=qt-unified-mac-x64-online
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$QT.dmg https://download.qt.io/official_releases/online_installers/$QT.dmg
  hdiutil attach -noverify -noautofsck -quiet /tmp/$QT.dmg
  QT=qt-unified-mac-x64-3.0.5-online
  if [ "$IOS" == "true" ] || [ "$IOS_SIMULATOR" == "true" ]
  then
    /Volumes/$QT/$QT.app/Contents/MacOS/$QT -v --script $GOPATH/src/github.com/peterq/pan-light/qt/internal/ci/iscript.qs IOS=true
  else
    /Volumes/$QT/$QT.app/Contents/MacOS/$QT -v --script $GOPATH/src/github.com/peterq/pan-light/qt/internal/ci/iscript.qs DARWIN=true
  fi
  diskutil unmountDisk disk1
  rm -f /tmp/$QT.dmg
  ln -s $HOME/Qt $HOME/Desktop
fi

if [ "$ANDROID" == "true" ]
then
  #download and install android sdk
  SDK=sdk-tools-darwin-3859397.zip
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$SDK https://dl.google.com/android/repository/$SDK
  unzip -qq /tmp/$SDK -d $HOME/android-sdk-macosx/
  rm -f /tmp/$SDK
  ln -s $HOME/android-sdk-macosx $HOME/Desktop

  #install deps for android sdk
  $HOME/android-sdk-macosx/tools/bin/sdkmanager --list --verbose
  echo "y" | $HOME/android-sdk-macosx/tools/bin/sdkmanager "platform-tools" "build-tools;26.0.0" "platforms;android-25"
  echo "y" | $HOME/android-sdk-macosx/tools/bin/sdkmanager --update

  #download and install android ndk
  NDK=android-ndk-r18b-darwin-x86_64.zip
  curl -sL --retry 10 --retry-delay 10 -o /tmp/$NDK https://dl.google.com/android/repository/$NDK
  unzip -qq /tmp/$NDK -d $HOME
  rm -f /tmp/$NDK
  ln -s $HOME/android-ndk-r18b $HOME/Desktop
fi

#prepare env
sudo chown $USER /usr/local/bin
sudo chown $USER $GOROOT/pkg | true

#check env
df -h
diskutil list

ls $HOME/*
du -sh $HOME/*

exit 0