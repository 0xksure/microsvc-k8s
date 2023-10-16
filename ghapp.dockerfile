FROM golang:1.21 AS build

WORKDIR /service
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build ./cmd/github-app/*.go

ENTRYPOINT ["./build"]