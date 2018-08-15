FROM golang:1.10 as build
WORKDIR /go/src/github.com/dpb587/boshua
COPY . .
ENV CGO_ENABLED=0
RUN go build -o /tmp/binaries/boshua ./main/boshua

FROM alpine:3.4 as binaries
RUN apk --no-cache add wget
RUN mkdir /tmp/binaries
RUN true \
  && wget -qO /tmp/binaries/bosh http://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-3.0.1-linux-amd64 \
  && echo "ccc893bab8b219e9e4a628ed044ebca6c6de9ca0  /tmp/binaries/bosh" | sha1sum -c \
  && chmod +x /tmp/binaries/bosh

FROM ubuntu:16.04
RUN true \
  && apt-get update \
  && apt-get install -y ca-certificates curl git openssh-client multipath-tools
COPY --from=binaries /tmp/binaries /usr/local/bin
COPY --from=build /tmp/binaries /usr/local/bin
