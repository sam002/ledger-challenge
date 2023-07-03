FROM golang:1.20-alpine AS build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/accounts cmd/accounts/accounts.go

FROM alpine:3.18
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=build /usr/local/bin/accounts ./
EXPOSE 3003
CMD ["./accounts"]