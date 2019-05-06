FROM golang AS builder

RUN mkdir /bluebell
WORKDIR /bluebell
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN go get -u github.com/gobuffalo/packr/v2/packr2
RUN packr2
RUN GOOS=linux GOARCH=amd64 go build -a -o /go/bin/bluebell

FROM golang:alpine
LABEL maintainer="mail@dongfg.com"

COPY --from=builder /go/bin/bluebell /go/bin/bluebell
ENTRYPOINT ["/go/bin/bluebell"]