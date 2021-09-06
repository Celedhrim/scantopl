FROM golang:alpine AS builder

WORKDIR $GOPATH/src/github.com/Celedhrim/scantopl/

COPY . .

RUN go get -d -v
RUN go build -o /go/bin/scantopl

FROM alpine
COPY --from=builder /go/bin/scantopl /usr/bin/scantopl

ARG UID=2001
ARG GID=2001
ARG UNAME=scanservjs

RUN addgroup -g $GID $UNAME
RUN adduser -u $UID -G $UNAME -h /output -D -s /bin/sh $UNAME


ENV \
  # The paperless instance URL
  PLURL="http://127.0.0.1:8080" \
  # The paperless token
  PLTOKEN="XXXXXXXXXXXXXXXXXXXXXXX"

USER $UNAME

ENTRYPOINT ["/usr/bin/scantopl", "-scandir", "/output"]
