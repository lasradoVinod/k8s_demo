FROM golang:1.20

WORKDIR /app

COPY helloworld/greeter_server/* ./

RUN go mod download
RUN go build -o main .

EXPOSE 50051

CMD ["/app/main"]
