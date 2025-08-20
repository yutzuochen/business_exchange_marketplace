#!/bin/bash

# 🔍 Cloud SQL 連接測試腳本

echo "🔍 測試 Cloud SQL 連接..."

# 檢查環境變數
if [ -z "$DB_HOST" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "❌ 錯誤: 請設置必要的資料庫環境變數"
    echo "   請設置: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME"
    echo ""
    echo "   例如:"
    echo "   export DB_HOST=YOUR_CLOUD_SQL_IP"
    echo "   export DB_USER=app"
    echo "   export DB_PASSWORD=your_password"
    echo "   export DB_NAME=business_exchange"
    exit 1
fi

echo "📋 連接資訊:"
echo "   主機: $DB_HOST"
echo "   用戶: $DB_USER"
echo "   資料庫: $DB_NAME"
echo "   埠號: ${DB_PORT:-3306}"

# 測試連接
echo ""
echo "🔌 測試資料庫連接..."

# 檢查是否安裝了 mysql 客戶端
if command -v mysql &> /dev/null; then
    echo "✅ 使用 mysql 客戶端測試連接..."
    
    # 測試連接
    if mysql -h "$DB_HOST" -P "${DB_PORT:-3306}" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SELECT 1 as test;" 2>/dev/null; then
        echo "✅ 資料庫連接成功!"
        
        # 顯示資料庫資訊
        echo ""
        echo "📊 資料庫資訊:"
        mysql -h "$DB_HOST" -P "${DB_PORT:-3306}" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "
        SELECT 
            'Database' as info_type,
            '$DB_NAME' as value
        UNION ALL
        SELECT 
            'Version' as info_type,
            VERSION() as value
        UNION ALL
        SELECT 
            'Current User' as info_type,
            USER() as value
        UNION ALL
        SELECT 
            'Current Database' as info_type,
            DATABASE() as value;
        " 2>/dev/null
        
        # 檢查表結構
        echo ""
        echo "📋 資料表列表:"
        mysql -h "$DB_HOST" -P "${DB_PORT:-3306}" -u "$DB_USER" -p"$DB_PASSWORD" "$DB_NAME" -e "SHOW TABLES;" 2>/dev/null
        
    else
        echo "❌ 資料庫連接失敗!"
        echo "   請檢查:"
        echo "   1. IP 地址是否正確"
        echo "   2. 用戶名和密碼是否正確"
        echo "   3. 資料庫是否存在"
        echo "   4. 網路連接是否正常"
        echo "   5. Cloud SQL 授權網路設置"
    fi
else
    echo "⚠️  未安裝 mysql 客戶端，無法直接測試連接"
    echo "   請安裝 MySQL 客戶端或使用其他工具測試"
    echo ""
    echo "   安裝方法:"
    echo "   Ubuntu/Debian: sudo apt install mysql-client"
    echo "   CentOS/RHEL: sudo yum install mysql"
    echo "   macOS: brew install mysql-client"
fi

echo ""
echo "🔧 其他測試方法:"
echo "   1. 使用 Cloud SQL Proxy:"
echo "      cloud_sql_proxy -instances=YOUR_INSTANCE_CONNECTION_NAME=tcp:3306"
echo ""
echo "   2. 使用 Adminer (如果可用):"
echo "      http://localhost:8081"
echo ""
echo "   3. 檢查應用程式日誌:"
echo "      docker compose logs app"
