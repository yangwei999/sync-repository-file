FROM golang:1.18.8 as BUILDER

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/sync-repository-file
COPY . .
RUN cd server/entrance && GO111MODULE=on CGO_ENABLED=0 go build -a -o sync-repository-file-server .

# copy binary config and utils
FROM alpine:3.14
COPY --from=BUILDER /go/src/github.com/opensourceways/sync-repository-file/server/entrance/sync-repository-file-server /opt/app/

ENTRYPOINT ["/opt/app/sync-repository-file-server"]