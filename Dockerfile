FROM golang:latest

LABEL org.opencontainers.image.source https://github.com/telenornms/skogul

RUN mkdir -p src

WORKDIR src

COPY go.mod go.sum ./

COPY . .

RUN make skogul

ENTRYPOINT ["./skogul"]
