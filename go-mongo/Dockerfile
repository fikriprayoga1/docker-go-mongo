FROM mongo

COPY shell.sh /usr/local
COPY go.mod /
COPY go.sum /

RUN apt-get update
RUN apt-get upgrade
RUN apt-get install curl
RUN curl -O https://dl.google.com/go/go1.16.3.linux-amd64.tar.gz
RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.3.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin
RUN chmod +x /usr/local/shell.sh
RUN /usr/local/go/bin/go version
RUN mkdir /usr/local/go/src/server

COPY server/. /usr/local/go/src/server

CMD /usr/local/shell.sh