FROM golang:1.26-alpine AS build-env
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server .

FROM alpine:3.21
WORKDIR /app
COPY --from=build-env /app/server .
CMD ["./server"]
