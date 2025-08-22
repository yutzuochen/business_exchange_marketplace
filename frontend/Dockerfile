# 使用官方 Node.js 18 作為基礎鏡像
FROM node:18-alpine AS base

# 安裝依賴
FROM base AS deps
RUN apk add --no-cache libc6-compat
WORKDIR /app

# 複製 package.json 和 package-lock.json
COPY package.json package-lock.json* ./
RUN npm ci

# 構建應用
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

# 設置環境變量
ENV NEXT_TELEMETRY_DISABLED=1
ENV NODE_ENV=production

# 設置 API URL 環境變量用於構建
ARG NEXT_PUBLIC_API_URL=https://business-exchange-backend-430730011391.us-central1.run.app
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

# 構建應用
RUN npm run build

# 生產環境
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# 設置運行時環境變數
ENV NEXT_PUBLIC_API_URL=https://business-exchange-backend-430730011391.us-central1.run.app

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# 複製構建產物
COPY --from=builder /app/public ./public

# 設置正確的權限
RUN mkdir .next
RUN chown nextjs:nodejs .next

# 複製構建產物
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
# v3