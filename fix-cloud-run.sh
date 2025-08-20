#!/bin/bash

# 🔧 修復 Cloud Run 服務腳本
# 專案 ID: businessexchange-468413

set -e

echo "🔧 開始修復 Cloud Run 服務..."

# 檢查登入狀態
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ 請先登入 Google Cloud:"
    echo "   gcloud auth login"
    exit 1
fi

# 設置專案
PROJECT_ID="businessexchange-468413"
REGION="us-central1"  # 使用實際的地區
SERVICE_NAME="trade-company"  # 使用實際的服務名稱

echo "📋 專案資訊:"
echo "   專案 ID: ${PROJECT_ID}"
echo "   地區: ${REGION}"
echo "   服務名稱: ${SERVICE_NAME}"

# 1. 檢查服務狀態
echo "📊 檢查服務狀態..."
gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format="table(metadata.name,status.url,status.conditions[0].status,status.conditions[0].message)"

# 2. 檢查服務日誌
echo "📝 檢查服務日誌..."
gcloud run services logs read ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --limit=50

# 3. 設置環境變數
echo "🔧 設置環境變數..."
gcloud run services update ${SERVICE_NAME} \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --set-env-vars "APP_ENV=production" \
    --set-env-vars "APP_NAME=BusinessExchange" \
    --set-env-vars "DB_HOST=10.80.0.3" \
    --set-env-vars "DB_USER=app" \
    --set-env-vars "DB_PASSWORD=app_password" \
    --set-env-vars "DB_NAME=business_exchange" \
    --set-env-vars "JWT_SECRET=your-production-secret-key" \
    --set-env-vars "REDIS_ADDR=10.80.0.3:6379"

# 4. 重新部署服務
echo "🚀 重新部署服務..."
gcloud run services update ${SERVICE_NAME} \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --memory 1Gi \
    --cpu 1 \
    --max-instances 10

# 5. 檢查修復結果
echo "✅ 修復完成!"
echo "🌐 服務 URL: $(gcloud run services describe ${SERVICE_NAME} --region ${REGION} --project ${PROJECT_ID} --format='value(status.url)')"

echo ""
echo "📝 下一步:"
echo "   1. 測試網站是否正常運作"
echo "   2. 檢查資料庫連接"
echo "   3. 查看服務日誌確認無錯誤"
