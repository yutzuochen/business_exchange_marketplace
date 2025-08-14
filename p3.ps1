啟動 Cloud SQL Proxy
cloud-sql-proxy "businessexchange-468413:us-central1:trade-sql" --port 3306

連線 MySQL
mysql -u app -p -h 127.0.0.1 -P 3306