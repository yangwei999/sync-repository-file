FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/sync-repository-file
COPY . .
RUN cd server/entrance && GO111MODULE=on CGO_ENABLED=0 go build -a -o sync-repository-file-server .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 sync-repository-file-server && \
    useradd -u 1000 -g sync-repository-file-server -s /bin/bash -m sync-repository-file-server

USER sync-repository-file-server

COPY --chown=sync-repository-file-server --from=BUILDER /go/src/github.com/opensourceways/sync-repository-file/server/entrance/sync-repository-file-server /opt/app/

ENTRYPOINT ["/opt/app/sync-repository-file-server"]