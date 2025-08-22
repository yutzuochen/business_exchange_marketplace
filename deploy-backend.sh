#!/bin/bash

# 設置變量
PROJECT_ID="businessexchange-468413"
REGION="us-central1"
SERVICE_NAME="business-exchange-backend"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

echo "🚀 開始部署後端到 Google Cloud..."
echo "專案 ID: ${PROJECT_ID}"
echo "地區: ${REGION}"
echo "服務名稱: ${SERVICE_NAME}"

# 檢查 gcloud 是否已登入
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" | grep -q .; then
    echo "❌ 請先登入 Google Cloud: gcloud auth login"
    exit 1
fi

# 設置專案
echo "📋 設置專案..."
gcloud config set project ${PROJECT_ID}

# 啟用必要的 API
echo "🔧 啟用必要的 API..."
gcloud services enable cloudbuild.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable containerregistry.googleapis.com

# 配置 Docker 認證
echo "🐳 配置 Docker 認證..."
gcloud auth configure-docker

echo ""
echo "=== 構建並推送後端 Docker 鏡像 ==="

# 構建 Docker 鏡像
echo "🔨 構建後端 Docker 鏡像..."
docker build -t ${IMAGE_NAME} .

if [ $? -eq 0 ]; then
    echo "✅ 後端 Docker 鏡像構建成功"
    
    # 推送 Docker 鏡像
    echo "📤 推送後端 Docker 鏡像..."
    docker push ${IMAGE_NAME}
    
    if [ $? -eq 0 ]; then
        echo "✅ 後端 Docker 鏡像推送成功"
        
        # 部署到 Cloud Run
        echo "🚀 部署後端到 Cloud Run..."
        gcloud run deploy ${SERVICE_NAME} \
            --image ${IMAGE_NAME} \
            --platform managed \
            --region ${REGION} \
            --allow-unauthenticated \
            --memory 1Gi \
            --cpu 1 \
            --max-instances 10 \
            --set-env-vars "APP_ENV=production,APP_NAME=BusinessExchange,DB_HOST=127.0.0.1,DB_PORT=3306,DB_USER=app,DB_PASSWORD=app_password,DB_NAME=business_exchange,CLOUDSQL_CONNECTION_NAME=${PROJECT_ID}:${REGION}-c:trade-sql" \
            --add-cloudsql-instances ${PROJECT_ID}:${REGION}-c:trade-sql \
            --timeout 300 \
            --cpu-boost
        
        if [ $? -eq 0 ]; then
            echo "✅ 後端部署成功！"
            BACKEND_URL=$(gcloud run services describe ${SERVICE_NAME} --region=${REGION} --format="value(status.url)")
            echo "🌐 後端 URL: ${BACKEND_URL}"
        else
            echo "❌ 後端部署失敗"
            echo "💡 檢查日誌: gcloud logging read 'resource.type=cloud_run_revision AND resource.labels.service_name=${SERVICE_NAME}' --limit=20"
        fi
    else
        echo "❌ 後端 Docker 鏡像推送失敗"
    fi
else
    echo "❌ 後端 Docker 鏡像構建失敗"
fi

echo ""
echo "🎉 部署腳本執行完成！"
