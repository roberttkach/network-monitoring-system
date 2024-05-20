FROM golang:latest

WORKDIR /app

COPY ./src ./src
COPY go.sum .
COPY go.mod .
COPY main.go .

RUN go build -o main .

CMD ["./main"]
