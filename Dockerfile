FROM golang:latest

WORKDIR /transaction-service

COPY . .

RUN go build ./cmd/transaction-service

CMD ["./transaction-service"]