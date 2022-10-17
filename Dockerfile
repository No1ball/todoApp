FROM golang:1.14-alpine

COPY . /app

WORKDIR /app

RUN go build -o server .

CMD ["/app/server"]