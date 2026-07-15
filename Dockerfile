FROM golang:1.26.5-alpine3.24 AS build

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app ./cmd/app


FROM alpine:3.24
COPY --from=build /app /app
EXPOSE 8080
CMD ["/app"]
