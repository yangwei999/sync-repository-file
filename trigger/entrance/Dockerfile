FROM openeuler/openeuler:23.03 as BUILDER
RUN dnf update -y && \
    dnf install -y golang && \
    go env -w GOPROXY=https://goproxy.cn,direct

MAINTAINER zengchen1024<chenzeng765@gmail.com>

# build binary
WORKDIR /go/src/github.com/opensourceways/sync-repository-file
COPY . .
RUN cd trigger/entrance && GO111MODULE=on CGO_ENABLED=0 go build -a -o sync-repository-file-trigger .

# copy binary config and utils
FROM openeuler/openeuler:22.03
RUN dnf -y update && \
    dnf in -y shadow && \
    groupadd -g 1000 sync-repository-file-trigger && \
    useradd -u 1000 -g sync-repository-file-trigger -s /bin/bash -m sync-repository-file-trigger

USER sync-repository-file-trigger

COPY --chown=sync-repository-file-trigger --from=BUILDER /go/src/github.com/opensourceways/sync-repository-file/trigger/entrance/sync-repository-file-trigger /opt/app/

ENTRYPOINT ["/opt/app/sync-repository-file-trigger"]