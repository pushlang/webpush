FROM golang:1.10

RUN mkdir -p /go/src/webpush

WORKDIR /go/src/webpush

COPY . /go/src/webpush

RUN go get -d -v ./...

RUN go install -v ./...

ENV DB_DSN "user=postgres password=12345 dbname=pushover sslmode=disable"

CMD ["~/go/bin/webpush -user=ucryge6j8mr9jnyhkef5jkab71y7sm"]

EXPOSE 8000
