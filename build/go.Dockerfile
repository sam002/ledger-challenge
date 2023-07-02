FROM golang:1.20-alpine
COPY go.mod go.sum ./
RUN go mod download && go mod verify