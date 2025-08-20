#!/bin/bash

# 🔍 測試 Cloud SQL 連接腳本
# 專案 ID: businessexchange-468413

set -e

echo "🔍 測試 Cloud SQL 連接..."

# 檢查必要的工具
if ! command -v mysql &> /dev/null; then
    echo "❌ MySQL 客戶端未安裝，嘗試安裝..."
    sudo apt update && sudo apt install -y mysql-client
fi

# 測試連接
echo "📡 測試連接到 Cloud SQL..."
echo "   主機: 10.80.0.3"
echo "   用戶: app"
echo "   資料庫: business_exchange"

# 測試連接
if mysql -h 10.80.0.3 -u app -papp_password -e "SELECT 1;" 2>/dev/null; then
    echo "✅ Cloud SQL 連接成功!"
    
    # 檢查資料庫
    echo "📊 檢查資料庫狀態..."
    mysql -h 10.80.0.3 -u app -papp_password -e "SHOW DATABASES;"
    
    # 檢查表
    echo "📋 檢查 business_exchange 資料庫中的表..."
    mysql -h 10.80.0.3 -u app -papp_password -e "USE business_exchange; SHOW TABLES;"
    
else
    echo "❌ Cloud SQL 連接失敗!"
    echo ""
    echo "🔧 可能的解決方案:"
    echo "   1. 檢查 Cloud SQL 實例是否運行"
    echo "   2. 檢查防火牆規則"
    echo "   3. 檢查 IP 地址是否正確"
    echo "   4. 檢查用戶權限"
    echo ""
    echo "📝 下一步:"
    echo "   1. 在 Google Cloud Console 中檢查 Cloud SQL 狀態"
    echo "   2. 確認實例的連接資訊"
    echo "   3. 檢查網路配置"
fi
