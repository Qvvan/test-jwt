FROM golang:1.23.1-alpine as builder

WORKDIR /usr/local/src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go build -o ./bin/app cmd/app/main.go

FROM alpine:latest

COPY --from=builder ./usr/local/src/bin/app /

CMD ["/app"]
