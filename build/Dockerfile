FROM rockylinux:9
ENV SHELL=bash
ENV PATH=$PATH:/usr/local/go/bin

RUN yum install -y vim man git python3-docutils rpm-build wget make

COPY go.mod /tmp/
RUN set -e; \
    GO_VERSION=$(awk '/^go / {print $2}' /tmp/go.mod); \
    wget -q https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz; \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz; \
    rm go${GO_VERSION}.linux-amd64.tar.gz

ADD . /skogul
