FROM golang:1.8.2
RUN go get github.com/laher/goxc
ENV USER root
WORKDIR /go/src/github.com/yuuki/albio
ADD . /go/src/github.com/yuuki/albio
