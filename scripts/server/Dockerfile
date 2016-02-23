FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install software-properties-common ca-certificates -y
RUN add-apt-repository ppa:ambrop7/badvpn -y
RUN apt-get update
RUN apt-get install -y badvpn curl

RUN apt-get install -y libappindicator3-1

ENV LANTERN_BINARY https://github.com/getlantern/lantern/releases/download/2.0.16/update_linux_amd64.bz2

RUN curl -L $LANTERN_BINARY | bzip2 -d - > /usr/bin/lantern
RUN chmod +x /usr/bin/lantern

COPY start.sh /usr/bin/start.sh

ENTRYPOINT ["/usr/bin/start.sh"]
