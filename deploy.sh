#!/bin/bash

# 設置變量
PROJECT_ID="businessexchange-468413"
REGION="us-central1"
BACKEND_SERVICE_NAME="business-exchange-backend"
FRONTEND_SERVICE_NAME="business-exchange-frontend"
BACKEND_IMAGE_NAME="gcr.io/${PROJECT_ID}/${BACKEND_SERVICE_NAME}"
FRONTEND_IMAGE_NAME="gcr.io/${PROJECT_ID}/${FRONTEND_SERVICE_NAME}"

echo "🚀 開始部署到 Google Cloud..."
echo "專案 ID: ${PROJECT_ID}"
echo "地區: ${REGION}"

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
echo "=== 部署後端 Go 專案 ==="

# 構建並推送後端 Docker 鏡像
echo "🔨 構建後端 Docker 鏡像..."
docker build -t ${BACKEND_IMAGE_NAME} .

if [ $? -eq 0 ]; then
    echo "✅ 後端 Docker 鏡像構建成功"
    
    echo "📤 推送後端 Docker 鏡像..."
    docker push ${BACKEND_IMAGE_NAME}
    
    if [ $? -eq 0 ]; then
        echo "✅ 後端 Docker 鏡像推送成功"
        
        # 部署到 Cloud Run
        echo "🚀 部署後端到 Cloud Run..."
        gcloud run deploy ${BACKEND_SERVICE_NAME} \
            --image ${BACKEND_IMAGE_NAME} \
            --platform managed \
            --region ${REGION} \
            --allow-unauthenticated \
            --memory 1Gi \
            --cpu 1 \
            --max-instances 10 \
            --set-env-vars "APP_ENV=production,APP_NAME=BusinessExchange,DB_HOST=127.0.0.1,DB_PORT=3306,DB_USER=app,DB_PASSWORD=app_password,DB_NAME=business_exchange,CLOUDSQL_CONNECTION_NAME=${PROJECT_ID}:${REGION}-c:trade-sql" \
            --add-cloudsql-instances ${PROJECT_ID}:${REGION}-c:trade-sql
        
        if [ $? -eq 0 ]; then
            echo "✅ 後端部署成功！"
            BACKEND_URL=$(gcloud run services describe ${BACKEND_SERVICE_NAME} --region=${REGION} --format="value(status.url)")
            echo "🌐 後端 URL: ${BACKEND_URL}"
        else
            echo "❌ 後端部署失敗"
        fi
    else
        echo "❌ 後端 Docker 鏡像推送失敗"
    fi
else
    echo "❌ 後端 Docker 鏡像構建失敗"
fi

echo ""
echo "=== 部署前端 Next.js 專案 ==="

# 構建並推送前端 Docker 鏡像
echo "🔨 構建前端 Docker 鏡像..."
cd frontend
docker build -t ${FRONTEND_IMAGE_NAME} .

if [ $? -eq 0 ]; then
    echo "✅ 前端 Docker 鏡像構建成功"
    
    echo "📤 推送前端 Docker 鏡像..."
    docker push ${FRONTEND_IMAGE_NAME}
    
    if [ $? -eq 0 ]; then
        echo "✅ 前端 Docker 鏡像推送成功"
        
        # 部署到 Cloud Run
        echo "🚀 部署前端到 Cloud Run..."
        gcloud run deploy ${FRONTEND_SERVICE_NAME} \
            --image ${FRONTEND_IMAGE_NAME} \
            --platform managed \
            --region ${REGION} \
            --allow-unauthenticated \
            --memory 1Gi \
            --cpu 1 \
            --max-instances 10 \
            --set-env-vars "NEXT_PUBLIC_API_URL=${BACKEND_URL:-http://localhost:8080}"
        
        if [ $? -eq 0 ]; then
            echo "✅ 前端部署成功！"
            FRONTEND_URL=$(gcloud run services describe ${FRONTEND_SERVICE_NAME} --region=${REGION} --format="value(status.url)")
            echo "🌐 前端 URL: ${FRONTEND_URL}"
        else
            echo "❌ 前端部署失敗"
        fi
    else
        echo "❌ 前端 Docker 鏡像推送失敗"
    fi
else
    echo "❌ 前端 Docker 鏡像構建失敗"
fi

cd ..

echo ""
echo "🎉 部署完成！"
echo "後端 URL: ${BACKEND_URL:-未獲取}"
echo "前端 URL: ${FRONTEND_URL:-未獲取}"
echo ""
echo "💡 提示："
echo "1. 確保 Cloud SQL 實例正在運行"
echo "2. 檢查環境變量是否正確設置"
echo "3. 測試 API 端點是否正常工作"
