'use client';

import { useState, useEffect } from 'react';

export default function DebugPage() {
  const [logs, setLogs] = useState<string[]>([]);
  
  const addLog = (message: string) => {
    console.log(message);
    setLogs(prev => [...prev, `${new Date().toLocaleTimeString()}: ${message}`]);
  };

  useEffect(() => {
    const runDiagnostic = async () => {
      addLog('🔍 開始診斷...');
      
      // 檢查環境變數
      addLog(`📝 NEXT_PUBLIC_API_URL: ${process.env.NEXT_PUBLIC_API_URL || '未設定'}`);
      
      // 檢查網路連接
      try {
        addLog('🌐 測試網路連接...');
        const response = await fetch('https://httpbin.org/get');
        addLog(`✅ 網路連接正常 (狀態: ${response.status})`);
      } catch (error) {
        addLog(`❌ 網路連接失敗: ${error}`);
      }
      
      // 測試後端 API
      try {
        addLog('🔗 測試後端 API...');
        const apiUrl = process.env.NEXT_PUBLIC_API_URL || 'https://business-exchange-backend-430730011391.us-central1.run.app';
        addLog(`📡 使用 API URL: ${apiUrl}`);
        
        const response = await fetch(`${apiUrl}/api/v1/listings`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        addLog(`📊 API 響應狀態: ${response.status}`);
        addLog(`📋 響應頭: ${JSON.stringify(Object.fromEntries(response.headers))}`);
        
        if (response.ok) {
          const data = await response.json();
          addLog(`✅ API 調用成功! 獲取到 ${data.listings?.length || 0} 個商品`);
        } else {
          addLog(`❌ API 調用失敗: ${response.status} ${response.statusText}`);
        }
      } catch (error) {
        addLog(`💥 API 調用異常: ${error instanceof Error ? error.message : String(error)}`);
      }
      
      addLog('🏁 診斷完成');
    };
    
    runDiagnostic();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8">
          <span className="text-orange-500">567</span> 診斷頁面
        </h1>
        
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-4">前端診斷日誌</h2>
          
          <div className="bg-black text-green-400 p-4 rounded font-mono text-sm max-h-96 overflow-y-auto">
            {logs.map((log, index) => (
              <div key={index} className="mb-1">
                {log}
              </div>
            ))}
            {logs.length === 0 && (
              <div className="text-gray-500">正在載入診斷信息...</div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
