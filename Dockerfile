FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/api/main.go

FROM alpine:3.14

ARG TLS_CERT
ARG TLS_KEY

WORKDIR /runtime
COPY $TLS_CERT $TLS_KEY ./certs/
COPY --from=build /app/main ./
EXPOSE ${PORT}
CMD ["./main"]