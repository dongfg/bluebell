FROM golang AS builder

RUN mkdir /bluebell
WORKDIR /bluebell
COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /go/bin/bluebell ./cmd/bluebell

FROM golang:alpine
LABEL maintainer="mail@dongfg.com"

COPY --from=builder /go/bin/bluebell /go/bin/bluebell
ENTRYPOINT ["/go/bin/bluebell", "/opt/config/config.yml"]