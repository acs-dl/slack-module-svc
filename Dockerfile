FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/acs-dl/slack-module-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/slack_client-module-svc /go/src/github.com/acs-dl/slack_client-module-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/slack-module-svc /usr/local/bin/slack-module-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["slack-module-svc"]
