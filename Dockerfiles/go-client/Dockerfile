FROM golang:1.20

WORKDIR /app

COPY helloworld/greeter_client/* ./

RUN go mod download
RUN go build -o main .

CMD ["/app/main", "-addr", "server-service.default.svc.cluster.local:50051"]
