FROM golang:latest

WORKDIR /app

VOLUME [ "./posts" ]

RUN go mod download
RUN go build github.com/Tch1b0/JP-backend

CMD ["./JP-backend"]