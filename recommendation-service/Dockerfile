FROM golang:1.23.1 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .
RUN go build -o main ./cmd/main.go

# smaller linux image
FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/main /
EXPOSE 8080
CMD ["/main"]