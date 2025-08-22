'use client';

import { useState, useEffect } from 'react';

export default function BasicPage() {
  const [count, setCount] = useState(0);
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    console.log('✅ BasicPage useEffect 執行了!');
    setMounted(true);
  }, []);

  const handleClick = () => {
    console.log('🖱️ 按鈕被點擊了!');
    setCount(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto text-center">
        <h1 className="text-4xl font-bold mb-8">
          <span className="text-orange-500">567</span> 基礎測試頁面
        </h1>
        
        <div className="bg-white rounded-lg shadow p-6 space-y-4">
          <p className="text-lg">
            組件狀態: {mounted ? '✅ 已掛載' : '⏳ 未掛載'}
          </p>
          
          <p className="text-lg">
            點擊次數: <span className="font-bold text-blue-600">{count}</span>
          </p>
          
          <button
            onClick={handleClick}
            className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors"
          >
            點擊測試 React 互動
          </button>
          
          <div className="mt-6 text-sm text-gray-600">
            <p>如果您看到這個頁面並且按鈕可以點擊，說明 React 正常工作</p>
            <p>請打開瀏覽器開發者工具的 Console 查看日誌</p>
          </div>
        </div>
      </div>
    </div>
  );
}


