FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build ./cmd/api/main.go

FROM alpine:3.14

WORKDIR /runtime
COPY --from=build /app/main ./
EXPOSE ${PORT}
CMD ["./main"]