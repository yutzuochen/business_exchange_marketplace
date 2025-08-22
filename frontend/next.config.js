/** @type {import('next').NextConfig} */
const nextConfig = {
  // 啟用獨立輸出，用於 Docker 部署
  output: 'standalone',
  
  // 禁用圖片優化（在 Cloud Run 中可能會有問題）
  images: {
    unoptimized: true,
  },
  
  // 環境變量
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
  },
  
  // 關閉 React Strict Mode 以避免水合問題
  reactStrictMode: false,
  
  // 關閉 SWC minify，使用 Terser
  swcMinify: false,
  
  // 優化構建配置
  compiler: {
    removeConsole: false,
  },
}

module.exports = nextConfig
