#!/bin/bash
set -ev

if [[ "$DESKTOP" == "true" ]]; then $GOPATH/bin/qtsetup full desktop; fi
if [[ "$ANDROID" == "true" ]]; then $GOPATH/bin/qtsetup full android; fi
if [[ "$IOS" == "true" ]]; then $GOPATH/bin/qtsetup full ios; fi
if [[ "$IOS_SIMULATOR" == "true" ]]; then $GOPATH/bin/qtsetup full ios-simulator; fi
if [[ "$QT_MXE" == "true" ]]; then $GOPATH/bin/qtsetup full windows; fi

exit 0
