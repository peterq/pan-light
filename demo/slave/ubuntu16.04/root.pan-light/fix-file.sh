#!/usr/bin/env bash
appname=`basename $0 | sed s,\.sh$,,`

dirname=`dirname $0`
tmp="${dirname#?}"

if [ "${dirname%$tmp}" != "/" ]; then
dirname=$PWD/$dirname
fi

cp ${dirname}/files/font-cn.ttf /usr/share/fonts/
cp ${dirname}/files/xstartup /root/.vnc/xstartup
