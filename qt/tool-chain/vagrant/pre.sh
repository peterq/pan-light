#!/bin/bash
set -ev

if [[ "$OS" == "darwin" ]]; then
  export PROF=.bash_profile
  export GO=go1.11.2.darwin-amd64.tar.gz
else if [[ "$OS" == "linux" ]]; then
  export PROF=.profile
  export GO=go1.11.2.linux-amd64.tar.gz

  sudo apt-get -qq update && sudo apt-get -y -qq install curl git software-properties-common libgl1-mesa-dev fontconfig unzip && sudo apt-get -qq clean

  if false; then
    sudo apt-get -qq update && sudo apt-get -y -qq install bison build-essential gperf flex ruby python libasound2-dev libbz2-dev libcap-dev libcups2-dev libdrm-dev libegl1-mesa-dev libgcrypt11-dev libnss3-dev libpci-dev libpulse-dev libudev-dev libxtst-dev gyp ninja-build && sudo apt-get -qq clean
    sudo apt-get -qq update && sudo apt-get -y -qq install libssl-dev libxcursor-dev libxcomposite-dev libxdamage-dev libxrandr-dev libfontconfig1-dev libxss-dev libsrtp0-dev libwebp-dev libjsoncpp-dev libopus-dev libavutil-dev libavformat-dev libavcodec-dev libevent-dev libxslt1-dev && sudo apt-get -qq clean

    sudo apt-get -qq update && sudo apt-get -y -qq install lxde xinit && sudo apt-get -qq clean
    sudo /usr/share/debconf/fix_db.pl #or sudo apt-get -y -qq remove miscfiles dictionaries-common
    echo "exec startlxde" >> $HOME/.xinitrc
    sudo startx &
  fi
fi; fi

#darwin
if [[ "$QT_HOMEBREW" == "true" ]]; then echo "export QT_HOMEBREW=true" >> $HOME/$PROF; fi

#linux
if [[ "$QT_PKG_CONFIG" == "true" ]]; then
  echo "export QT_PKG_CONFIG=true" >> $HOME/$PROF
  echo "export PKG_CONFIG_PATH=/opt/qt510/lib/pkgconfig" >> $HOME/$PROF
  echo "export QT_DOC_DIR=/opt/qt510/doc" >> $HOME/$PROF
  echo "export QT_MISC_DIR=/opt/qt510" >> $HOME/$PROF
fi

if [[ "$QT_MXE" == "true" ]]; then
  echo "export QT_MXE_ARCH="$QT_MXE_ARCH >> $HOME/$PROF
  echo "export QT_MXE_STATIC="$QT_MXE_STATIC >> $HOME/$PROF
fi

curl -sL --retry 10 --retry-delay 10 -o /tmp/$GO https://dl.google.com/go/$GO && tar -xzf /tmp/$GO -C $HOME && rm -f /tmp/$GO

echo "export PATH=$PATH:$HOME/go/bin" >> $HOME/$PROF
echo "export GOROOT=$HOME/go" >> $HOME/$PROF
echo "export GOPATH=$HOME/gopath" >> $HOME/$PROF
source $HOME/$PROF

if [[ "$OS" == "darwin" ]]; then
  ln -s $HOME/go $HOME/Desktop/GOROOT
  ln -s $HOME/gopath $HOME/Desktop/GOPATH
fi

go get -v -tags=no_env github.com/peterq/pan-light/qt/cmd/...

if [[ "$OS" == "darwin" ]]; then
  sudo xcode-select -s /Applications/Xcode.app/Contents/Developer
  if [[ "$IOS" == "true" ]]; then rm -R /Applications/Xcode.app/Contents/Developer/Platforms/AppleTVOS.platform; rm -R /Applications/Xcode.app/Contents/Developer/Platforms/WatchOS.platform; fi
  if [[ "$ANDROID" == "true" ]]; then brew update; brew tap caskroom/versions; brew cask install java8; fi
else if [[ "$OS" == "linux" ]]; then
  sudo rm -f -R $HOME/.config
  sudo rm -f -R $HOME/.cache
fi; fi

$GOPATH/bin/qtsetup prep

exit 0
