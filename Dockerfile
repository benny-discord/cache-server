FROM golang:1.15
RUN mkdir /app
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o main ./...
CMD ["/app/main"]
