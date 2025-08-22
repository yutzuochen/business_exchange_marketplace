'use client';

import { useState, useEffect } from 'react';

export default function SimpleMarketPage() {
  const [status, setStatus] = useState('開始載入...');
  const [data, setData] = useState<{listings?: Array<{title: string, price: number, category: string}>, pagination?: {page: number, total_pages: number}} | null>(null);

  useEffect(() => {
    const testAPI = async () => {
      try {
        setStatus('正在測試API連接...');
        
        const apiUrl = 'https://business-exchange-backend-430730011391.us-central1.run.app/api/v1/listings';
        console.log('Calling API:', apiUrl);
        
        const response = await fetch(apiUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        console.log('Response status:', response.status);
        console.log('Response headers:', response.headers);
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const result = await response.json();
        console.log('API response:', result);
        
        setData(result);
        setStatus(`成功載入 ${result.listings?.length || 0} 個商品`);
        
      } catch (error) {
        console.error('API call failed:', error);
        setStatus(`載入失敗: ${error instanceof Error ? error.message : '未知錯誤'}`);
      }
    };

    testAPI();
  }, []);

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <h1 className="text-4xl font-bold text-center mb-8">
          <span className="text-orange-500">567</span> 企業互惠平台
        </h1>
        
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-2xl font-semibold mb-4">市場測試頁面</h2>
          
          <div className="mb-4">
            <strong>狀態:</strong> {status}
          </div>
          
          {data && (
            <div className="mt-6">
              <h3 className="text-lg font-semibold mb-2">API 響應數據:</h3>
              <div className="bg-gray-100 p-4 rounded text-sm overflow-auto max-h-96">
                <p><strong>總商品數:</strong> {data.listings?.length || 0}</p>
                <p><strong>分頁信息:</strong> 第 {data.pagination?.page || 'N/A'} 頁，共 {data.pagination?.total_pages || 'N/A'} 頁</p>
                
                {data.listings && data.listings.length > 0 && (
                  <div className="mt-4">
                    <strong>商品列表:</strong>
                    <ul className="mt-2 space-y-1">
                      {data.listings.slice(0, 5).map((listing, index: number) => (
                        <li key={index} className="text-gray-700">
                          • {listing.title} - ${(listing.price / 100).toFixed(2)} ({listing.category})
                        </li>
                      ))}
                      {data.listings.length > 5 && (
                        <li className="text-gray-500">... 還有 {data.listings.length - 5} 個商品</li>
                      )}
                    </ul>
                  </div>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
