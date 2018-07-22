FROM golang:latest

WORKDIR /go/src/app
COPY . .

RUN go get -v github.com/weaming/proxyany

EXPOSE 20443

CMD ["proxyany"]
