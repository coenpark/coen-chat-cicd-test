# syntax=docker/dockerfile:1

# Build the application from source
FROM --platform=linux/amd64 golang:1.20 AS build-stage
ENV GO111MODULE=on

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /chat .

# Run the tests in the container
FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM atlassian/ubuntu-minimal:latest AS build-release-stage
WORKDIR /src
COPY src/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /chat .

# Run the tests in the container
FROM build-stage AS run-test-stage
#RUN go test -v ./...

# Deploy the application binary into a lean image
FROM atlassian/ubuntu-minimal:latest AS build-release-stage

WORKDIR /

COPY --from=build-stage /chat /chat

EXPOSE 8080

ENTRYPOINT ["/chat"]