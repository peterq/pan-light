#!/bin/bash
appname=`basename $0 | sed s,\.sh$,,`

dirname=`dirname $0`
tmp="${dirname#?}"

if [ "${dirname%$tmp}" != "/" ]; then
dirname=$PWD/$dirname
fi
export LD_LIBRARY_PATH="$dirname/deploy/linux/lib"
export QT_PLUGIN_PATH="$dirname/deploy/linux/plugins"
export QML_IMPORT_PATH="$dirname/deploy/linux/qml"
export QML2_IMPORT_PATH="$dirname/deploy/linux/qml"
export DISPLAY=:1

cd /root
/usr/bin/vncserver :1 -geometry 1600x900 -depth 24

cd $dirname
$dirname/$appname "$@"
