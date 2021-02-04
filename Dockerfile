#
# 1. Build Container
#
FROM golang:1.15.7 AS build

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# First add modules list to better utilize caching
COPY go.sum go.mod /app/

WORKDIR /app

# Download dependencies
RUN go mod download

COPY . /app

# Build components.
# Put built binaries and runtime resources in /app dir ready to be copied over or used.
RUN go install -race -installsuffix cgo -ldflags="-w -s" && \
    mkdir -p /app && \
    cp -r $GOPATH/bin/api-tutorme /app/

#
# 2. Runtime Container
#
FROM alpine

ENV PATH="/app:${PATH}"

# See http://stackoverflow.com/questions/34729748/installed-go-binary-not-found-in-path-on-alpine-linux-docker
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

WORKDIR /app

COPY --from=build /app /app/

EXPOSE 8010

CMD ["./api-tutorme"]