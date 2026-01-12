ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}

LABEL org.opencontainers.image.source https://github.com/telenornms/skogul

RUN mkdir -p src

WORKDIR src

COPY go.mod go.sum ./

COPY . .

RUN make skogul

ENTRYPOINT ["./skogul"]
