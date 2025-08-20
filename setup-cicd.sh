#!/bin/bash

# 🚀 CI/CD 快速設置腳本

set -e

echo "🚀 開始設置 CI/CD 流程..."

# 檢查必要的工具
check_tools() {
    echo "🔍 檢查必要工具..."
    
    if ! command -v git &> /dev/null; then
        echo "❌ Git 未安裝"
        exit 1
    fi
    
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker 未安裝"
        exit 1
    fi
    
    echo "✅ 所有必要工具已安裝"
}

# 設置 Git 分支
setup_branches() {
    echo "🌿 設置 Git 分支..."
    
    # 檢查當前分支
    CURRENT_BRANCH=$(git branch --show-current)
    echo "   當前分支: $CURRENT_BRANCH"
    
    # 創建 develop 分支
    if ! git branch | grep -q "develop"; then
        echo "   創建 develop 分支..."
        git checkout -b develop
        git push -u origin develop
    else
        echo "   develop 分支已存在"
    fi
    
    # 返回原分支
    git checkout $CURRENT_BRANCH
    
    echo "✅ 分支設置完成"
}

# 設置 GitHub Secrets 提示
setup_secrets_info() {
    echo ""
    echo "🔐 需要在 GitHub Repository 中設置以下 Secrets:"
    echo ""
    echo "   1. GCP_SA_KEY: Google Cloud 服務帳戶 JSON 金鑰"
    echo "   2. DB_HOST: Cloud SQL 實例 IP"
    echo "   3. DB_USER: 資料庫用戶名"
    echo "   4. DB_PASSWORD: 資料庫密碼"
    echo "   5. DB_NAME: 資料庫名稱"
    echo "   6. JWT_SECRET: JWT 密鑰"
    echo "   7. REDIS_ADDR: Redis 地址"
    echo "   8. SLACK_WEBHOOK: Slack Webhook URL (可選)"
    echo ""
    echo "   設置路徑: Settings > Secrets and variables > Actions"
}

# 設置環境保護規則
setup_environments() {
    echo ""
    echo "🌍 需要在 GitHub Repository 中設置環境:"
    echo ""
    echo "   1. staging: 測試環境"
    echo "   2. production: 生產環境"
    echo ""
    echo "   設置路徑: Settings > Environments"
}

# 創建服務帳戶
create_service_account() {
    echo ""
    echo "🔑 創建 Google Cloud 服務帳戶..."
    echo ""
    echo "   執行以下命令創建服務帳戶:"
    echo ""
    echo "   # 創建服務帳戶"
    echo "   gcloud iam service-accounts create business-exchange-sa \\"
    echo "       --display-name='Business Exchange Service Account'"
    echo ""
    echo "   # 設置權限"
    echo "   gcloud projects add-iam-policy-binding businessexchange-468413 \\"
    echo "       --member='serviceAccount:business-exchange-sa@businessexchange-468413.iam.gserviceaccount.com' \\"
    echo "       --role='roles/run.admin'"
    echo ""
    echo "   gcloud projects add-iam-policy-binding businessexchange-468413 \\"
    echo "       --member='serviceAccount:business-exchange-sa@businessexchange-468413.iam.gserviceaccount.com' \\"
    echo "       --role='roles/storage.admin'"
    echo ""
    echo "   gcloud projects add-iam-policy-binding businessexchange-468413 \\"
    echo "       --member='serviceAccount:business-exchange-sa@businessexchange-468413.iam.gserviceaccount.com' \\"
    echo "       --role='roles/iam.serviceAccountUser'"
    echo ""
    echo "   # 創建金鑰"
    echo "   gcloud iam service-accounts keys create ~/business-exchange-sa.json \\"
    echo "       --iam-account=business-exchange-sa@businessexchange-468413.iam.gserviceaccount.com"
    echo ""
    echo "   然後將 ~/business-exchange-sa.json 的內容複製到 GitHub Secrets 的 GCP_SA_KEY"
}

# 測試本地構建
test_local_build() {
    echo ""
    echo "🧪 測試本地構建..."
    
    if docker build -t business-exchange:test .; then
        echo "✅ 本地構建成功"
    else
        echo "❌ 本地構建失敗"
        exit 1
    fi
}

# 顯示下一步
show_next_steps() {
    echo ""
    echo "🎯 下一步操作:"
    echo ""
    echo "   1. 在 GitHub Repository 中設置 Secrets"
    echo "   2. 設置環境保護規則"
    echo "   3. 創建 Google Cloud 服務帳戶"
    echo "   4. 推送代碼到 develop 分支測試 CI/CD"
    echo "   5. 合併到 main 分支部署到生產環境"
    echo ""
    echo "📚 詳細文檔: CI_CD_SETUP.md"
    echo "🚀 部署腳本: deploy-to-cloud.sh"
}

# 主函數
main() {
    echo "🚀 BusinessExchange CI/CD 設置"
    echo "================================"
    
    check_tools
    setup_branches
    setup_secrets_info
    setup_environments
    create_service_account
    test_local_build
    show_next_steps
    
    echo ""
    echo "✅ CI/CD 設置完成!"
}

# 執行主函數
main
