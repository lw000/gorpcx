FROM golang:1.17

WORKDIR /app

ADD  . /app

RUN go mod tidy

RUN go build -o gorpcx ./main.go

EXPOSE 8001

CMD ["/app/gorpcx"]