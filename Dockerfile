FROM golang:1.24.4-alpine3.22

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o /main

COPY *.go ./
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]