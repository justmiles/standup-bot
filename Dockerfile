FROM golang:1.12-stretch as builder

COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -installsuffix cgo .
RUN md5sum standup-bot

# Create image from scratch
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/standup-bot /standup-bot
COPY --from=builder /tmp /tmp

ENV AWS_DEFAULT_REGION us-east-1

ENTRYPOINT ["/standup-bot"]
