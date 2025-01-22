FROM golang:1.23.5 as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN env GOOS=linux CGO_ENABLED=0 go build -o musicApp ./cmd
RUN chmod +x /app/musicApp


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/musicApp /app
COPY --from=builder /app/migration/schemas /app/migration/schemas

CMD [ "/app/musicApp" ]
