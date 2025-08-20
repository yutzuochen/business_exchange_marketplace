#!/bin/bash

# 🔧 修復 Cloud Run 環境變數腳本
# 專案 ID: businessexchange-468413

set -e

echo "🔧 開始修復 Cloud Run 環境變數..."

# 檢查登入狀態
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ 請先登入 Google Cloud:"
    echo "   gcloud auth login"
    exit 1
fi

# 設置專案
PROJECT_ID="businessexchange-468413"
REGION="us-central1"
SERVICE_NAME="trade-company"

echo "📋 專案資訊:"
echo "   專案 ID: ${PROJECT_ID}"
echo "   地區: ${REGION}"
echo "   服務名稱: ${SERVICE_NAME}"

# 檢查當前環境變數
echo "📊 當前環境變數:"
gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format="value(spec.template.spec.containers[0].env[].name,spec.template.spec.containers[0].env[].value)"

# 創建環境變數文件
echo "📝 創建環境變數文件..."
cat > env-vars.yaml << 'EOF'
APP_ENV: "production"
APP_NAME: "BusinessExchange"
APP_PORT: "8080"
DB_HOST: "10.80.0.3"
DB_PORT: "3306"
DB_USER: "app"
DB_PASSWORD: "app_password"
DB_NAME: "business_exchange"
DB_SSL_MODE: "require"
REDIS_ADDR: "10.80.0.3:6379"
REDIS_PASSWORD: ""
REDIS_DB: "0"
JWT_SECRET: "your-production-secret-key-change-this"
JWT_ISSUER: "businessexchange-468413"
JWT_EXPIRY: "24h"
CORS_ALLOW_ORIGINS: "*"
CORS_ALLOW_METHODS: "GET,POST,PUT,DELETE,OPTIONS"
CORS_ALLOW_HEADERS: "*"
FILE_UPLOAD_MAX_SIZE: "10MB"
FILE_UPLOAD_ALLOWED_TYPES: "image/jpeg,image/png,image/gif"
PAGINATION_DEFAULT_LIMIT: "20"
PAGINATION_MAX_LIMIT: "100"
EOF

# 設置環境變數
echo "🔧 設置環境變數..."
gcloud run services update ${SERVICE_NAME} \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --env-vars-file env-vars.yaml

# 清理臨時文件
rm -f env-vars.yaml

# 檢查修復結果
echo "✅ 環境變數設置完成!"
echo "📊 更新後的環境變數:"
gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format="value(spec.template.spec.containers[0].env[].name,spec.template.spec.containers[0].env[].value)"

echo ""
echo "🌐 服務 URL: $(gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format='value(status.url)')"

echo ""
echo "📝 下一步:"
echo "   1. 等待服務重新啟動 (約 1-2 分鐘)"
echo "   2. 測試網站是否正常運作"
echo "   3. 檢查服務日誌確認無錯誤"
echo "   4. 如果還有問題，檢查 Cloud SQL 連接"
