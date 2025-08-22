'use client';

import { useState, useEffect } from 'react';
import { Listing } from '@/types/listing';
import { apiClient } from '@/lib/api';
import ListingCard from '@/components/ListingCard';
import Link from 'next/link';

export default function MarketPage() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [categories, setCategories] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      console.log('🚀 Starting data fetch...');
      setLoading(true);
      setError(null);
      
      try {
        console.log('📡 Calling API...');
        
        // 使用 Promise.allSettled 確保即使部分失敗也能繼續
        const [listingsResult, categoriesResult] = await Promise.allSettled([
          apiClient.getListings(),
          apiClient.getCategories(),
        ]);
        
        console.log('📊 Results:', { listingsResult, categoriesResult });
        
        // 處理 listings 結果
        if (listingsResult.status === 'fulfilled') {
          setListings(listingsResult.value || []);
          console.log('✅ Listings loaded:', listingsResult.value?.length || 0);
        } else {
          console.error('❌ Listings failed:', listingsResult.reason);
        }
        
        // 處理 categories 結果
        if (categoriesResult.status === 'fulfilled') {
          setCategories(categoriesResult.value || []);
          console.log('✅ Categories loaded:', categoriesResult.value?.length || 0);
        } else {
          console.error('❌ Categories failed:', categoriesResult.reason);
        }
        
        // 如果兩個都失敗了才設置錯誤
        if (listingsResult.status === 'rejected' && categoriesResult.status === 'rejected') {
          setError('無法載入數據');
        }
        
      } catch (err) {
        console.error('💥 Unexpected error:', err);
        setError(err instanceof Error ? err.message : '載入失敗');
      } finally {
        console.log('🏁 Setting loading to false');
        setLoading(false);
      }
    };

    // 添加延遲確保組件完全掛載
    const timer = setTimeout(fetchData, 100);
    return () => clearTimeout(timer);
  }, []);

  const filteredListings = selectedCategory === 'all' 
    ? listings 
    : listings.filter(listing => listing.category === selectedCategory);

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">載入中...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-600 text-xl mb-4">載入失敗</div>
          <p className="text-gray-600 mb-4">{error}</p>
          <div className="text-sm text-gray-500 mb-4">
            <p>請檢查：</p>
            <ul className="list-disc list-inside mt-2">
              <li>後端服務是否運行在 http://192.168.168.169:8080</li>
              <li>API 端點是否正常響應</li>
              <li>網絡連接是否正常</li>
            </ul>
          </div>
          <button 
            onClick={() => window.location.reload()} 
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            重新載入
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      {/* 頁面標題 */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-4">
              <h1 className="text-6xl font-bold text-orange-500 drop-shadow-lg">567</h1>
              <span className="text-2xl font-bold text-orange-500">我來接</span>
            </div>
            <div className="text-center">
              <p className="text-lg font-semibold text-gray-800">企業互惠平台</p>
            </div>
            <Link 
              href="/"
              className="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700"
            >
              返回首頁
            </Link>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        {/* 分類篩選 */}
        <div className="mb-8">
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => setSelectedCategory('all')}
              className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                selectedCategory === 'all'
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
              }`}
            >
              全部
            </button>
            {Array.isArray(categories) && categories.map((category) => (
              <button
                key={category}
                onClick={() => setSelectedCategory(category)}
                className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
                  selectedCategory === category
                    ? 'bg-blue-600 text-white'
                    : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
                }`}
              >
                {category}
              </button>
            ))}
          </div>
        </div>

        {/* 結果統計 */}
        <div className="mb-6">
          <p className="text-gray-600">
            顯示 {filteredListings.length} 個結果
            {selectedCategory !== 'all' && ` (${selectedCategory} 分類)`}
          </p>
        </div>

        {/* Listings 網格 */}
        {filteredListings.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredListings.map((listing) => (
              <ListingCard key={listing.id} listing={listing} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12">
            <div className="text-gray-400 text-6xl mb-4">📭</div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">沒有找到結果</h3>
            <p className="text-gray-600">
              {selectedCategory === 'all' 
                ? '目前沒有任何 listings' 
                : `在 ${selectedCategory} 分類中沒有找到結果`}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

