FROM golang:1.18 as builder

WORKDIR /app

COPY ${PWD} .

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o loggen main.go

FROM golang:1.18-alpine

EXPOSE 8080

WORKDIR /
COPY --from=builder /app/loggen .

ENTRYPOINT ["./loggen"]
