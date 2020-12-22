FROM golang:1.15-alpine as build

ENV GOARCH=amd64 \
    GOPROXY=goproxy.cn,direct

WORKDIR /go/src

COPY go.sum .
COPY go.mod .
RUN go mod download

COPY . .
RUN go build -o gproxy main.go

FROM docker:20

WORKDIR /go/bin

COPY --chown=0:0 --from=build /go/src/gproxy /go/bin/goproxy
ENTRYPOINT ["./goproxy"]

CMD ["http", "--port", "80"]
