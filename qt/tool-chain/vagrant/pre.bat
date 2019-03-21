
::enable delayed expansion
setlocal enabledelayedexpansion


::disable updates
reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows\CurrentVersion\WindowsUpdate\Auto Update" /v AUOptions /t REG_DWORD /d 1 /f
sc config wuauserv start= disabled


::install curl
expand -f:* c:\tmp\curl.cab c:\tmp\
mv -f /cygdrive/c/tmp/AMD64/* "/cygdrive/c/Program Files/OpenSSH/bin/"


::install 7z
set SZ=7z1805-x64.exe
curl -sL --retry 10 --retry-delay 10 -o %TMP%\%SZ% https://7-zip.org/a/%SZ%
%TMP%\%SZ% /S
del %TMP%\%SZ% /Q
setx /M PATH "%PATH%;C:\Progra~1\7-Zip"
set PATH=%PATH%;C:\Progra~1\7-Zip


::install Git
set GIT=Git-2.19.0-64-bit.exe
curl -sL --retry 10 --retry-delay 10 -o %TMP%\%GIT% https://github.com/git-for-windows/git/releases/download/v2.19.0.windows.1/%GIT%
%TMP%\%GIT% /silent /norestart
del %TMP%\%GIT% /Q
setx /M PATH "%PATH%;C:\Progra~1\Git\bin"
set PATH=%PATH%;C:\Progra~1\Git\bin


::install Go + pull repo
set GO=go1.11.2.windows-amd64.msi
curl -sL --retry 10 --retry-delay 10 -o %TMP%\%GO% https://storage.googleapis.com/golang/%GO%
%TMP%\%GO% /passive /norestart
del %TMP%\%GO% /Q
reg delete "HKEY_CURRENT_USER\Environment" /v GOPATH /f
setx /M PATH "%PATH%;C:\Go\bin"
set PATH=%PATH%;C:\Go\bin
setx /M GOPATH "C:\gopath"
set GOPATH=C:\gopath
setx /M GOROOT "C:\go"
set GOROOT=C:\go

go get -v -tags=no_env github.com/peterq/pan-light/qt/cmd/...


::install VC++ 2015 Redis
set VC=vc_redist.x64.exe
curl -sL --retry 10 --retry-delay 10 -o %TMP%\%VC% https://download.microsoft.com/download/9/3/F/93FCF1E7-E6A4-478B-96E7-D4B285925B00/%VC%
%TMP%\%VC% /passive /norestart
del %TMP%\%VC% /Q


if "%QT_MSYS2%" == "true" (
  setx /M QT_MSYS2 "%QT_MSYS2%"
  setx /M QT_MSYS2_STATIC "%QT_MSYS2_STATIC%"
  setx /M QT_MSYS2_ARCH "%QT_MSYS2_ARCH%"

  if "%QT_MSYS2_ARCH%" == "386" (
    setx /M MSYSTEM "MINGW32"
    echo MSYSTEM=MINGW32>> C:\Users\vagrant\.ssh\environment
  ) else (
    setx /M MSYSTEM "MINGW64"
    echo MSYSTEM=MINGW64>> C:\Users\vagrant\.ssh\environment
  )

  echo QT_MSYS2=true>> C:\Users\vagrant\.ssh\environment
  echo QT_MSYS2_STATIC=%QT_MSYS2_STATIC%>> C:\Users\vagrant\.ssh\environment
  echo QT_MSYS2_ARCH=%QT_MSYS2_ARCH%>> C:\Users\vagrant\.ssh\environment


  ::install msys2
  set MSYS2=msys2-x86_64-20180531.exe
  set AI=auto-install.js
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!MSYS2! http://repo.msys2.org/distrib/x86_64/!MSYS2!
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!AI! https://raw.githubusercontent.com/msys2/msys2-installer/master/!AI!
  %TMP%\!MSYS2! --script %TMP%\!AI!
  del %TMP%\!MSYS2! /Q
  del %TMP%\!AI! /Q


  C:\msys64\usr\bin\bash -l -c "pacman -Q"
  C:\msys64\usr\bin\bash -l -c "pacman -Syyu --noconfirm --noprogress"
  C:\msys64\usr\bin\bash -l -c "pacman -Syyu --noconfirm --noprogress"

  if "%QT_MSYS2_ARCH%" == "386" (
    if "%QT_MSYS2_STATIC%" == "true" (
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-i686-qt-creator mingw-w64-i686-qt5-static"
    ) else (
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-i686-qt-creator mingw-w64-i686-qt5"
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-i686-qtwebkit"
    )
  ) else (
    if "%QT_MSYS2_STATIC%" == "true" (
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-x86_64-qt-creator mingw-w64-x86_64-qt5-static"
    ) else (
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-x86_64-qt-creator mingw-w64-x86_64-qt5"
      C:\msys64\usr\bin\bash -l -c "pacman -S --noconfirm --needed --noprogress mingw-w64-x86-qtwebkit"
    )
  )

  C:\msys64\usr\bin\bash -l -c "pacman -Q"
  C:\msys64\usr\bin\bash -l -c "pacman -Scc --noconfirm --noprogress"
) else (
  ::install Qt
  set QT=qt-unified-windows-x86-online.exe
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!QT! https://download.qt.io/official_releases/online_installers/!QT!
  %TMP%\!QT! -v --script %GOPATH%\src\github.com\therecipe\qt\internal\ci\iscript.qs WINDOWS=true
  del %TMP%\!QT! /Q
  setx /M PATH "%PATH%;C:\Qt\Tools\mingw730_64\bin"
  set PATH=%PATH%;C:\Qt\Tools\mingw730_64\bin
)


::update ssh env variables
echo TMP=C:/tmp>> C:\Users\vagrant\.ssh\environment
net stop "OpenSSH Server"
net start "OpenSSH Server"


if "%ANDROID%" == "true" (
  ::install JDK
  set JDK=jdk-8u192-ea-bin-b04-windows-x64-01_aug_2018.exe
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!JDK! https://download.java.net/java/jdk8u192/archive/b04/binaries/!JDK!
  %TMP%\!JDK! /s
  del %TMP%\!JDK! /Q

  setx /M JAVA_HOME "C:\Progra~1\Java\jdk1.8.0_192"
  set JAVA_HOME=C:\Progra~1\Java\jdk1.8.0_192

  ::install Android SDK
  set SDK=sdk-tools-windows-4333796.zip
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!SDK! https://dl.google.com/android/repository/!SDK!
  7z x %TMP%\!SDK! -oC:\android-sdk-windows\
  del %TMP%\!SDK! /Q

  mkdir C:\android-sdk-windows\licenses
  echo fc946e8f231f3e3159bf0b7c655c924cb2e38330>> C:\android-sdk-windows\licenses\android-googletv-license
  echo d56f5187479451eabf01fb78af6dfcb131a6481e>> C:\android-sdk-windows\licenses\android-sdk-license
  echo 504667f4c0de7af1a06de9f4b1727b84351f2910>> C:\android-sdk-windows\licenses\android-sdk-preview-license
  echo 33b6a2b64607f11b759f320ef9dff4ae5c47d97a>> C:\android-sdk-windows\licenses\google-gdk-license
  echo d975f751698a77b662f1254ddbeed3901e976f5a>> C:\android-sdk-windows\licenses\intel-android-extra-license
  echo 63d703f5692fd891d5acacfbd8e09f40fc976105>> C:\android-sdk-windows\licenses\mips-android-sysimage-license

  cmd /C "C:\android-sdk-windows\tools\bin\sdkmanager.bat --list --verbose"
  cmd /C "C:\android-sdk-windows\tools\bin\sdkmanager.bat "platform-tools" "build-tools;26.0.0" "platforms;android-25""
  cmd /C "mv C:\android-sdk-windows\tools\ C:\android-sdk-windows\toolsOLD\"
  cmd /C "C:\android-sdk-windows\toolsOLD\bin\sdkmanager.bat --update"
  cmd /C "rm -R C:\android-sdk-windows\toolsOLD\"


  ::install Android NDK
  set NDK=android-ndk-r18b-windows-x86_64.zip
  curl -sL --retry 10 --retry-delay 10 -o %TMP%\!NDK! https://dl.google.com/android/repository/!NDK!
  7z x %TMP%\!NDK! -oC:\
  del %TMP%\!NDK! /Q
)


::qtsetup
if "%QT_MSYS2%" == "true" (
  if "%QT_MSYS2_ARCH%" == "386" (
    set MSYSTEM=MINGW32
    C:\msys64\usr\bin\bash -l -c "$GOPATH/bin/qtsetup full desktop"
  ) else (
    set MSYSTEM=MINGW64
    C:\msys64\usr\bin\bash -l -c "$GOPATH/bin/qtsetup full desktop"
  )
) else (
  if "%DESKTOP%" == "true" "%GOPATH%\bin\qtsetup" full desktop
  if "%ANDROID%" == "true" "%GOPATH%\bin\qtsetup" full android
)
