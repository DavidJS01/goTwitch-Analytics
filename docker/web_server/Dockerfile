# syntax=docker/dockerfile:1

# use golang alpine image
FROM golang:1.18-alpine AS builder

# COPY scripts/setup_grafana.sh .

# # install bash
# RUN apk add --no-cache bash
# # run setup scripts
# RUN ./setup_grafana.sh

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY . .
RUN go mod download

# Copy the code into the container.

# # Set necessary environment variables needed 
# # for our image and build the server.
# ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build ./go-twitch-web/api/main.go
RUN ls

# # # Command to run when starting the container.
ENTRYPOINT ["./main"]
