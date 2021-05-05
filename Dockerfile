FROM golang

WORKDIR /go/src

COPY server/. ./

RUN chmod 777 /go/src/shell.sh

CMD /go/src/shell.sh