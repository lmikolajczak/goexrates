FROM golang:1.20

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 /usr/bin/migrate

# Will install "column" which is make's dependency
RUN apt-get update && apt-get --yes install bsdmainutils

WORKDIR /code
COPY . .

RUN go get ./...
RUN cd cmd/api && go build -o /go/bin/goexrates-api
RUN cd cmd/cli && go build -o /go/bin/goexrates-cli

CMD tail -f /dev/null
