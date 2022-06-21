FROM golang:latest

RUN mkdir -p src

WORKDIR src

COPY go.mod go.sum ./

COPY . .

RUN make skogul

ENTRYPOINT ["./skogul"]
