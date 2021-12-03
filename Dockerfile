FROM golang:1.17

WORKDIR /src
ADD . .

RUN go build -o /main main.go
CMD ["/main"]