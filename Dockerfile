FROM golang:latest

WORKDIR /transaction-service

COPY . .

RUN go build ./cmd/transaction-service

EXPOSE 8080

CMD ["./transaction-service"]