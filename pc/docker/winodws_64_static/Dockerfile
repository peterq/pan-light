FROM therecipe/qt:windows_64_static
RUN groupadd -r -g 1000 peterq
RUN useradd -r -u 1000 -g peterq peterq
#USER peterq
#ENV HOME=/home/peterq
ENV GOPROXY=https://goproxy.io
VOLUME ["/home/user"]
WORKDIR /media/pan-light/pc
ENTRYPOINT ["/bin/bash"]
