FROM docker:cli AS cli
FROM golang:latest

COPY --from=cli /usr/local/bin/docker /usr/local/bin/docker

WORKDIR /etc/v2k/platform/sm
