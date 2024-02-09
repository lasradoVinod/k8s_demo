FROM golang:1.20

WORKDIR /app

COPY helloworld/greeter_server/* ./

RUN go mod download
RUN go build -o main .

EXPOSE 50051

CMD ["/app/main", "-target", "server-service.default.svc.cluster.local:50051", \
    "-Cert" , "/etc/secret/node-client-cert.pem" , "-Key", \
    "/etc/secret/node-client-key.pem", "-caKey", "/etc/secret/ca.pem"]
