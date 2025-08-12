$PROJECT="businessexchange-468413"
$REGION="us-central1"
$REPO="trade-repo"
$IMAGE="trade_company"
$TAG=$(git rev-parse --short HEAD)
$REGISTRY="${REGION}-docker.pkg.dev"

# Cloud SQL (MySQL, private IP)
gcloud sql instances create trade-sql --database-version=MYSQL_8_0 --region=$REGION --cpu=1 --memory=4GiB --network=default --no-assign-ip
gcloud sql databases create trade_company --instance=trade-sql
gcloud sql users create app --instance=trade-sql --password=app_password
$DB_IP=$(gcloud sql instances describe trade-sql --format="get(ipAddresses[?type='PRIVATE'].ipAddress)")

# Memorystore (Redis)
gcloud redis instances create trade-redis --region=$REGION --tier=basic --size=1 --network=default
$REDIS_HOST=$(gcloud redis instances describe trade-redis --region=$REGION --format="get(host)")

# Serverless VPC connector
gcloud compute networks vpc-access connectors create cr-connector --region=$REGION --network=default --range=10.8.0.0/28

# Deploy to Cloud Run (use the image you already pushed)
$REGISTRY="${REGION}-docker.pkg.dev"
gcloud run deploy trade-company `
  --image="${REGISTRY}/${PROJECT}/${REPO}/${IMAGE}:${TAG}" `
  --platform=managed --allow-unauthenticated `
  --vpc-connector=cr-connector --egress-settings=all-traffic `
  --set-env-vars="APP_ENV=production,APP_PORT=8080,DB_HOST=${DB_IP},DB_PORT=3306,DB_USER=app,DB_PASSWORD=app_password,DB_NAME=trade_company,REDIS_ADDR=${REDIS_HOST}:6379,REDIS_DB=0,JWT_SECRET=556611,CORS_ALLOWED_ORIGINS=*"

# Artifact Registry + build/push
gcloud artifacts repositories create $REPO --repository-format=docker --location=$REGION --quiet
$REGISTRY="${REGION}-docker.pkg.dev"
gcloud auth configure-docker $REGISTRY

docker build -t "${REGISTRY}/${PROJECT}/${REPO}/${IMAGE}:${TAG}" .
docker push "${REGISTRY}/${PROJECT}/${REPO}/${IMAGE}:${TAG}"


