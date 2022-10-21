# syntax=docker/dockerfile:1

# use golang alpine image
FROM golang:1.18-alpine

# create workdirectory in container, default path for commands
WORKDIR /app 

# copy needed Go files into /app
COPY cmd/api/main ./

EXPOSE 9090

