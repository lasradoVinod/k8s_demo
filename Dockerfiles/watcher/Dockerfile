FROM golang:1.20

WORKDIR /app

COPY watcher/* ./

RUN go mod download
RUN go build -o main .

CMD ["/app/main"]

#CMD ["tail", "-f", "/dev/null"]
