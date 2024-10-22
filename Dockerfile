FROM golang:1.23 AS builder

ENV DB_USER=
ENV DB_PASSWORD=
ENV DB_ADDR=

WORKDIR /app

COPY go.mod go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o crawler ./cmd

FROM gcr.io/distroless/static-debian11

WORKDIR /
COPY --from=builder /app/crawler /crawler

CMD ["/crawler"]






