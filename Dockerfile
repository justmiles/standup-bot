FROM golang:1.16-alpine

WORKDIR /app

COPY lib lib
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY LICENSE LICENSE

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /standup-bot

CMD [ "/standup-bot" ]
