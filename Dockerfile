FROM alpine:latest
WORKDIR /app
RUN mkdir /app/download
COPY wyze-go /app/wyze-go
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add --no-cache tzdata
ENV TZ=America/New_York
ENV WYZE_HOME=/app/download/
CMD /app/wyze-go
