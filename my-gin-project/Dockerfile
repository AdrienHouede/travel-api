FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src ./src

WORKDIR /app/src

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]