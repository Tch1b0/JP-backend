FROM golang:latest

WORKDIR /app

COPY ./main.go ./main.go
COPY ./account.json ./account.json
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
COPY ./pkg ./pkg
COPY ./posts ./posts

VOLUME [ "/app/posts" ]

RUN go mod download
RUN go build github.com/Tch1b0/JP-backend

CMD ["./JP-backend"]