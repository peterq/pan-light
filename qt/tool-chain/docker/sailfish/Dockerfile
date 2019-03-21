FROM i386/ubuntu:14.04 as base

ENV USER user
ENV HOME /home/$USER
ENV GOPATH $HOME/work

RUN apt-get -qq update && apt-get --no-install-recommends -qq -y install ca-certificates curl git
RUN GO=go1.11.2.linux-386.tar.gz && curl -sL --retry 10 --retry-delay 60 -O https://dl.google.com/go/$GO && tar -xzf $GO -C /usr/local
RUN /usr/local/go/bin/go get -tags=no_env github.com/peterq/pan-light/qt/cmd/...


FROM i386/ubuntu:14.04
LABEL maintainer therecipe

ENV USER user
ENV HOME /home/$USER
ENV GOPATH $HOME/work
ENV SF_SDK_TOOLING_DIR /srv/mer/toolings/SailfishOS-2.2.1.18
ENV PATH /usr/local/go/bin:$PATH
ENV QT_DOCKER true
ENV QT_SAILFISH true

COPY --from=base /usr/local/go /usr/local/go
COPY --from=base $GOPATH/bin $GOPATH/bin
COPY --from=base $GOPATH/src/github.com/peterq/pan-light/qt $GOPATH/src/github.com/peterq/pan-light/qt

RUN apt-get -qq update && apt-get --no-install-recommends -qq -y install ca-certificates curl && apt-get -qq clean
RUN apt-get -qq update && apt-get -y -qq purge python && apt-get -qq clean

RUN SF_SDK_TOOLING=Jolla-latest-Sailfish_SDK_Tooling-i486.tar.bz2 && mkdir -p $SF_SDK_TOOLING_DIR && curl -sL --retry 10 --retry-delay 60 -O https://releases.sailfishos.org/sdk/latest/$SF_SDK_TOOLING && tar --numeric-owner -p -xjf $SF_SDK_TOOLING -C $SF_SDK_TOOLING_DIR && rm -f $SF_SDK_TOOLING

RUN SF_SDK_TARGET=Jolla-latest-Sailfish_SDK_Target-i486.tar.bz2 && mkdir -p /srv/mer/targets/SailfishOS-2.2.1.18-i486 && curl -sL --retry 10 --retry-delay 60 -O https://releases.sailfishos.org/sdk/latest/$SF_SDK_TARGET && tar --numeric-owner -p -xjf $SF_SDK_TARGET -C /srv/mer/targets/SailfishOS-2.2.1.18-i486 && rm -f $SF_SDK_TARGET

RUN SF_SDK_TARGET=Jolla-latest-Sailfish_SDK_Target-armv7hl.tar.bz2 && mkdir -p /srv/mer/targets/SailfishOS-2.2.1.18-armv7hl && curl -sL --retry 10 --retry-delay 60 -O https://releases.sailfishos.org/sdk/latest/$SF_SDK_TARGET && tar --numeric-owner -p -xjf $SF_SDK_TARGET -C /srv/mer/targets/SailfishOS-2.2.1.18-armv7hl && rm -f $SF_SDK_TARGET

ENV PATH $SF_SDK_TOOLING_DIR/usr/bin/:$PATH

RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libmpc.so.3 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libmpfr.so.4 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libgmp.so.10 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libpthread_nonshared.a /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libc_nonshared.a /usr/lib/

RUN ln -s $SF_SDK_TOOLING_DIR/opt/cross/bin/i486-meego-linux-gnu-as $SF_SDK_TOOLING_DIR/opt/cross/libexec/gcc/i486-meego-linux-gnu/4.8.3/as
RUN ln -s $SF_SDK_TOOLING_DIR/opt/cross/bin/i486-meego-linux-gnu-ld $SF_SDK_TOOLING_DIR/opt/cross/libexec/gcc/i486-meego-linux-gnu/4.8.3/ld

RUN ln -s $SF_SDK_TOOLING_DIR/opt/cross/bin/armv7hl-meego-linux-gnueabi-as $SF_SDK_TOOLING_DIR/opt/cross/libexec/gcc/armv7hl-meego-linux-gnueabi/4.8.3/as
RUN ln -s $SF_SDK_TOOLING_DIR/opt/cross/bin/armv7hl-meego-linux-gnueabi-ld $SF_SDK_TOOLING_DIR/opt/cross/libexec/gcc/armv7hl-meego-linux-gnueabi/4.8.3/ld

RUN cd /srv/mer/targets/SailfishOS-2.2.1.18-i486/ && $SF_SDK_TOOLING_DIR/usr/bin/sb2-init -L "--sysroot=/" -C "--sysroot=/" -n -N -t / i486-meego-linux-gnu $SF_SDK_TOOLING_DIR/opt/cross/bin/i486-meego-linux-gnu-gcc

RUN cd /srv/mer/targets/SailfishOS-2.2.1.18-i486/ && $SF_SDK_TOOLING_DIR/usr/bin/sb2-init -L "--sysroot=/" -C "--sysroot=/" -n -N -t / armv7hl-meego-linux $SF_SDK_TOOLING_DIR/opt/cross/bin/i486-meego-linux-gnu-gcc

RUN sed -i 's/--target=$build_tgt/--target=$OPT_TARGET/g' /srv/mer/toolings/SailfishOS-2.2.1.18/usr/bin/mb2

RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/librpmbuild.so.8 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/librpm.so.8 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/librpmio.so.8 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libdb-4.8.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/liblua-5.1.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libpython2.7.so.1.0 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libbfd-2.25.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libnss3.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libelf.so.1 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libnssutil3.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libplc4.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libplds4.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libnspr4.so /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libglib-2.0.so.0 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libpcre.so.1 /usr/lib/
RUN ln -s $SF_SDK_TOOLING_DIR/usr/lib/libdw.so.1 /usr/lib/

RUN $GOPATH/bin/qtsetup prep

RUN $GOPATH/bin/qtsetup generate sailfish
RUN $GOPATH/bin/qtsetup install sailfish
RUN cd $GOPATH/src/github.com/peterq/pan-light/qt/internal/examples/sailfish/listview && $GOPATH/bin/qtdeploy build sailfish && rm -rf ./deploy

RUN $GOPATH/bin/qtsetup generate sailfish-emulator
RUN $GOPATH/bin/qtsetup install sailfish-emulator
RUN cd $GOPATH/src/github.com/peterq/pan-light/qt/internal/examples/sailfish/listview && $GOPATH/bin/qtdeploy build sailfish-emulator && rm -rf ./deploy
