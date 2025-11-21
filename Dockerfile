FROM golang:alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/emil-server

EXPOSE 8888

CMD ["./main"]