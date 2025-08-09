# Build stage
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download || true
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/server /app/server
COPY --from=builder /app/templates /app/templates
ENV APP_ENV=production
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/server"] 