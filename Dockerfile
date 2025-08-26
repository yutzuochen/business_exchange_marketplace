# syntax=docker/dockerfile:1.7

########################
# 1) Build stage
########################
FROM golang:1.23-alpine AS builder
WORKDIR /src

# 基本工具
RUN apk add --no-cache git ca-certificates tzdata

# 先複製模組檔以利用快取
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 再複製其餘程式碼
COPY . .

# 版本資訊（可選）
ARG VERSION=dev
ARG COMMIT=none
ARG DATE=unknown

# 建置（依你的 main 調整：常見為 ./cmd/server）
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}" \
      -o /out/server ./cmd/server

########################
# 2) Runtime stage
########################
FROM alpine:3.20 AS runner
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# 非 root 使用者
RUN addgroup -S app && adduser -S -G app -u 10001 app
COPY --from=builder /out/server /app/server

ENV APP_ENV=production \
    GIN_MODE=release \
    PORT=8080 \
    TZ=Asia/Taipei

USER app
EXPOSE 8080
CMD ["/app/server"]
