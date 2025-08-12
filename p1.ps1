$PROJECT="businessexchange-468413"
$REGION="us-central1"
$REPO="trade-repo"
$IMAGE="trade_company"
$TAG=$(git rev-parse --short HEAD)

gcloud auth login
gcloud config set project $PROJECT
gcloud config set run/region $REGION
gcloud services enable run.googleapis.com artifactregistry.googleapis.com cloudbuild.googleapis.com sqladmin.googleapis.com redis.googleapis.com vpcaccess.googleapis.com servicenetworking.googleapis.com
gcloud artifacts repositories create $REPO --repository-format=docker --location=$REGION --quiet

$REGISTRY="${REGION}-docker.pkg.dev"
gcloud auth configure-docker $REGISTRY

docker build -t "${REGISTRY}/${PROJECT}/${REPO}/${IMAGE}:${TAG}" .
docker push "${REGISTRY}/${PROJECT}/${REPO}/${IMAGE}:${TAG}"