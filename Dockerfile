FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build server.go
RUN mkdir /app

CMD ["/go/src/app/server"]