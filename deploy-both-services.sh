#!/bin/bash

# 🚀 Deploy both backend and frontend to Google Cloud Run
set -e

PROJECT_ID="businessexchange-468413"
REGION="us-central1"

echo "🚀 Deploying BusinessExchange Backend and Frontend to Google Cloud..."

# 1. Deploy Backend First
echo "🔨 Building and deploying backend..."
docker build -t gcr.io/${PROJECT_ID}/business-exchange-backend .
docker push gcr.io/${PROJECT_ID}/business-exchange-backend

gcloud run deploy business-exchange-backend \
    --image gcr.io/${PROJECT_ID}/business-exchange-backend \
    --platform managed \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --allow-unauthenticated \
    --memory 1Gi \
    --cpu 1 \
    --max-instances 10 \
    --set-env-vars "APP_ENV=production,APP_NAME=BusinessExchange,DB_HOST=127.0.0.1,DB_PORT=3306,DB_USER=app,DB_PASSWORD=app_password,DB_NAME=business_exchange,CLOUDSQL_CONNECTION_NAME=businessexchange-468413:us-central1-c:trade-sql" \
    --add-cloudsql-instances businessexchange-468413:us-central1-c:trade-sql

# Get backend URL
BACKEND_URL=$(gcloud run services describe business-exchange-backend --region ${REGION} --project ${PROJECT_ID} --format="value(status.url)")
echo "✅ Backend deployed at: ${BACKEND_URL}"

# 2. Deploy Frontend with correct API URL
echo "🔨 Building and deploying frontend..."
cd frontend

# Build frontend Docker image with backend URL
docker build -t gcr.io/${PROJECT_ID}/business-exchange-frontend .
docker push gcr.io/${PROJECT_ID}/business-exchange-frontend

gcloud run deploy business-exchange-frontend \
    --image gcr.io/${PROJECT_ID}/business-exchange-frontend \
    --platform managed \
    --region ${REGION} \
    --project ${PROJECT_ID} \
    --allow-unauthenticated \
    --memory 512Mi \
    --cpu 1 \
    --max-instances 10 \
    --set-env-vars "NEXT_PUBLIC_API_URL=${BACKEND_URL}"

# Get frontend URL
FRONTEND_URL=$(gcloud run services describe business-exchange-frontend --region ${REGION} --project ${PROJECT_ID} --format="value(status.url)")

cd ..

echo "✅ Deployment completed!"
echo "🌐 Backend URL: ${BACKEND_URL}"
echo "🌐 Frontend URL: ${FRONTEND_URL}"
echo ""
echo "🧪 Testing services..."
echo "Backend health: $(curl -s ${BACKEND_URL}/api/v1/listings | head -c 50)..."
echo "Frontend status: $(curl -s ${FRONTEND_URL} | head -c 50)..."
