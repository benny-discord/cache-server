FROM golang:1.15
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./...
CMD ["/app/main"]
