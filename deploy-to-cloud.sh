#!/bin/bash

# 🚀 BusinessExchange 部署到 Google Cloud 腳本
# 專案 ID: businessexchange-468413

set -e

echo "🚀 開始部署 BusinessExchange 到 Google Cloud..."

# 檢查必要的環境變數
if [ -z "$GOOGLE_APPLICATION_CREDENTIALS" ]; then
    echo "❌ 錯誤: 請設置 GOOGLE_APPLICATION_CREDENTIALS 環境變數"
    echo "   例如: export GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account-key.json"
    exit 1
fi

# 設置專案
PROJECT_ID="businessexchange-468413"
REGION="us-central1"
SERVICE_NAME="business-exchange-backend"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"
CLOUDSQL_CONNECTION_NAME="${PROJECT_ID}:${REGION}:trade-sql"

echo "📋 專案資訊:"
echo "   專案 ID: ${PROJECT_ID}"
echo "   地區: ${REGION}"
echo "   服務名稱: ${SERVICE_NAME}"
echo "   映像名稱: ${IMAGE_NAME}"
echo "   Cloud SQL 連接: ${CLOUDSQL_CONNECTION_NAME}"

# 1. 構建 Docker 映像
echo "🔨 構建 Docker 映像..."
docker build -t ${IMAGE_NAME} .

# 2. 推送到 Google Container Registry
echo "📤 推送映像到 Google Container Registry..."
docker push ${IMAGE_NAME}

# 3. 部署到 Cloud Run
echo "🚀 部署到 Cloud Run..."
gcloud run deploy ${SERVICE_NAME} \
    --image ${IMAGE_NAME} \
    --platform managed \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --allow-unauthenticated \
    --memory 1Gi \
    --cpu 1 \
    --max-instances 10 \
    --set-env-vars "APP_ENV=production" \
    --set-env-vars "APP_NAME=BusinessExchange" \
    --set-env-vars "DB_HOST=/cloudsql/${CLOUDSQL_CONNECTION_NAME}" \
    --set-env-vars "DB_USER=app" \
    --set-env-vars "DB_PASSWORD=app_password" \
    --set-env-vars "DB_NAME=business_exchange" \
    --set-env-vars "JWT_ISSUER=${PROJECT_ID}" \
    --set-env-vars "JWT_SECRET=your-production-jwt-secret-change-me" \
    --set-env-vars "CORS_ALLOWED_ORIGINS=*" \
    --add-cloudsql-instances ${CLOUDSQL_CONNECTION_NAME}

# 4. 獲取服務 URL
SERVICE_URL=$(gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format="value(status.url)")

echo "✅ 部署完成!"
echo "🌐 服務 URL: ${SERVICE_URL}"
echo ""
echo "📝 重要配置說明:"
echo "   ✅ 服務器綁定: 0.0.0.0:\$PORT (Cloud Run 自動設置)"
echo "   ✅ 數據庫連接: Unix Socket ${CLOUDSQL_CONNECTION_NAME}"
echo "   ✅ Cloud SQL Proxy: 自動注入"
echo ""
echo "📝 下一步:"
echo "   1. 測試服務: curl ${SERVICE_URL}/health"
echo "   2. 檢查日誌: gcloud logs read --service=${SERVICE_NAME}"
echo "   3. 設置自定義域名 (可選)"

# 5. 顯示服務狀態
echo ""
echo "📊 服務狀態:"
gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format="table(metadata.name,status.url,status.conditions[0].status,status.conditions[0].message)"

# 6. 測試健康檢查
echo ""
echo "🔍 測試健康檢查..."
sleep 10
if curl -f "${SERVICE_URL}/health" 2>/dev/null; then
    echo "✅ 健康檢查通過"
else
    echo "⚠️  健康檢查失敗，請檢查日誌"
    echo "   查看日誌: gcloud logs read --service=${SERVICE_NAME} --limit=50"
fi
