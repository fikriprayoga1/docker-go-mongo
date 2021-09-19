FROM golang

WORKDIR /go/src

COPY src/. ./

CMD ["go", "run", "/go/src/profile.go"]