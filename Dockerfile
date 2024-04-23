FROM golang:1.21.4-alpine3.18 AS build
WORKDIR /build
RUN mkdir /build/bin
COPY . .
ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH
RUN go build -ldflags "-s" -o /build/bin/wyze-go
FROM alpine:latest
WORKDIR /app
RUN mkdir /app/download
COPY --from=build /build/bin/wyze-go /app/wyze-go
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache tzdata
ENV TZ=America/New_York
ENV WYZE_HOME=/app/download/
CMD /app/wyze-go
