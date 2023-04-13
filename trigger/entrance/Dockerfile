FROM golang:1.18.8 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/sync-repository-file
COPY . .
RUN cd trigger/entrance && GO111MODULE=on CGO_ENABLED=0 go build -a -o sync-repository-file-trigger .

# copy binary config and utils
FROM alpine:3.14
COPY --from=BUILDER /go/src/github.com/opensourceways/sync-repository-file/trigger/entrance/sync-repository-file-trigger /opt/app/

ENTRYPOINT ["/opt/app/sync-repository-file-trigger"]