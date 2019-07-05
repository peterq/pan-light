
FROM queeno/ubuntu-desktop
MAINTAINER Peter Q <me@peterq.cn>

RUN sed -i 's#http://archive.ubuntu.com/#http://cn.archive.ubuntu.com/#' /etc/apt/sources.list;

RUN apt-get update

RUN apt-get install -y libqt5multimedia5-plugins

WORKDIR /root/pan-light

COPY root.pan-light/ /root/pan-light/

RUN sh /root/pan-light/fix-file.sh

CMD /root/pan-light/demo_instance_manager.sh
