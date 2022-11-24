FROM ghcr.io/telenornms/skogul:latest
RUN mkdir dist
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o skogul ./cmd/skogul

FROM alpine:latest

#Install ca-certficates and dumb-init.
RUN apk --no-cache add ca-certificates dumb-init
RUN mkdir -p /etc/skogul/conf.d/

#Create symlink for musl.
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /

COPY ./docs/examples/basics/default_ipv4.json /etc/skogul/conf.d/default.json
COPY --from=0 /go/src/skogul /usr/local/bin/
RUN chmod +x /usr/local/bin/skogul

RUN addgroup -g 450 -S skogul && adduser -s /bin/sh -SD -G skogul skogul

ENTRYPOINT ["/usr/bin/dumb-init"]
CMD ["/usr/local/bin/skogul"]
