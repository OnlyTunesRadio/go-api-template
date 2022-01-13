FROM golang:1.17 as builder

RUN apt update && \
    apt install -y --no-install-recommends \
    ca-certificates

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod tidy

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o go-api-template ./main.go
RUN mv go-api-template go-api-template.bin

FROM ubuntu:20.04

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/go-api-template.bin /srv/go-api-template
COPY --from=builder /app/static/swagger.json /srv/static/swagger.json

EXPOSE 4000

CMD ["/srv/go-api-template"]