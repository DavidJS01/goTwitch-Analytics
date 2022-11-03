# ./Dockerfile-consumer

FROM golang:1.16-alpine AS builder

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed 
# for our image and build the consumer.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o main cmd/rabbitmq/consumer/main.go

RUN go build -o stream cmd/twitch/livestream/main.go

# Command to run when starting the container.
ENTRYPOINT ["./main"]
