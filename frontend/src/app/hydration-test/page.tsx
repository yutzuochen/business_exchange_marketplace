'use client';

import { useState, useEffect } from 'react';

export default function HydrationTest() {
  const [hydrated, setHydrated] = useState(false);
  const [clientTime, setClientTime] = useState<string>('未設置');

  useEffect(() => {
    console.log('🔄 HydrationTest useEffect 開始執行');
    
    // 設置水合狀態
    setHydrated(true);
    
    // 設置客戶端時間
    setClientTime(new Date().toLocaleTimeString());
    
    console.log('✅ HydrationTest useEffect 執行完成');
    
    // 測試定時器
    const timer = setInterval(() => {
      console.log('⏰ Timer tick');
      setClientTime(new Date().toLocaleTimeString());
    }, 1000);
    
    return () => {
      console.log('🧹 清理定時器');
      clearInterval(timer);
    };
  }, []);

  console.log('🎨 HydrationTest 渲染中，hydrated:', hydrated);

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-2xl mx-auto bg-white rounded-lg shadow p-6">
        <h1 className="text-3xl font-bold mb-6 text-center">
          <span className="text-blue-600">水合測試頁面</span>
        </h1>
        
        <div className="space-y-4 text-lg">
          <div className="p-4 border rounded">
            <strong>水合狀態: </strong>
            <span className={hydrated ? 'text-green-600' : 'text-red-600'}>
              {hydrated ? '✅ 已水合' : '❌ 未水合'}
            </span>
          </div>
          
          <div className="p-4 border rounded">
            <strong>客戶端時間: </strong>
            <span className="text-blue-600">{clientTime}</span>
          </div>
          
          <div className="p-4 border rounded">
            <strong>JavaScript 狀態: </strong>
            <span className="text-green-600">
              {typeof window !== 'undefined' ? '✅ 客戶端環境' : '❌ 服務端環境'}
            </span>
          </div>
          
          <div className="p-4 border rounded">
            <strong>渲染環境: </strong>
            <span className="text-purple-600">
              {typeof window !== 'undefined' ? '瀏覽器' : 'Node.js'}
            </span>
          </div>
        </div>
        
        <div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded">
          <p className="text-sm text-gray-600">
            <strong>調試說明：</strong><br/>
            1. 如果看到"未水合"，說明 useEffect 沒有執行<br/>
            2. 如果時間不更新，說明定時器沒有工作<br/>
            3. 請打開瀏覽器開發者工具查看 Console 日誌
          </p>
        </div>
      </div>
    </div>
  );
}
