#!/bin/bash

# Check if tables exist in Cloud SQL
echo "🔍 Checking database tables..."

# Connect to Cloud SQL and list tables
gcloud sql connect trade-sql --user=app --project=businessexchange-468413 --database=business_exchange << 'EOF'
SHOW TABLES;
EXIT
EOF
