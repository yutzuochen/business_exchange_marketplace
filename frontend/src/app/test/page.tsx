'use client';

import { useState } from 'react';
import { apiClient } from '@/lib/api';

export default function TestPage() {
  const [testResult, setTestResult] = useState<string>('');
  const [loading, setLoading] = useState(false);

  const testAPI = async () => {
    setLoading(true);
    setTestResult('測試中...');
    
    try {
      // 測試 listings API
      const listingsData = await apiClient.getListings();
      console.log('Listings API 成功:', listingsData);
      
      // 測試 categories API
      const categoriesData = await apiClient.getCategories();
      console.log('Categories API 成功:', categoriesData);
      
      setTestResult(`✅ API 測試成功！
      
Listings: ${listingsData.length} 個
Categories: ${categoriesData.length} 個

詳細信息請查看瀏覽器控制台。`);
    } catch (error) {
      console.error('API 測試失敗:', error);
      setTestResult(`❌ API 測試失敗: ${error instanceof Error ? error.message : '未知錯誤'}`);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 頁面標題 */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center justify-center space-x-4">
            <h1 className="text-6xl font-bold text-orange-500 drop-shadow-lg">567</h1>
            <span className="text-2xl font-bold text-orange-500">我來接</span>
          </div>
          <p className="mt-4 text-center text-lg font-semibold text-gray-800">企業互惠平台</p>
        </div>
      </div>

      <div className="max-w-4xl mx-auto px-4 py-16">
        <div className="text-center">
          <h2 className="text-3xl font-bold text-gray-900 mb-6">
            API 連接測試
          </h2>
          <p className="text-lg text-gray-600 mb-8">
            測試前端與後端 API 的連接狀態
          </p>

          <button
            onClick={testAPI}
            disabled={loading}
            className="bg-blue-600 text-white px-8 py-4 rounded-lg text-lg font-semibold hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? '測試中...' : '開始測試'}
          </button>

          {testResult && (
            <div className="mt-8 p-6 bg-white rounded-lg shadow-md">
              <pre className="whitespace-pre-wrap text-sm text-gray-800 font-mono">
                {testResult}
              </pre>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
