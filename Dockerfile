# Start from the latest golang base image
FROM golang:1.18.5

ENV GOPATH=/


## Copy go mod and sum files
COPY ./ ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download


# Build the Go app
RUN go build -o ./abf3 ./cmd/main.go

CMD ["./abf3"]
