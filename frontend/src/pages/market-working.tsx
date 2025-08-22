import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Listing } from '@/types/listing';
import { apiClient } from '@/lib/api';
import ListingCard from '@/components/ListingCard';

export default function MarketWorking() {
  const [listings, setListings] = useState<Listing[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedCategory, setSelectedCategory] = useState<string>('all');
  const [categories, setCategories] = useState<string[]>([]);

  useEffect(() => {
    const fetchData = async () => {
      console.log('🚀 MarketWorking: Starting data fetch...');
      setLoading(true);
      setError(null);
      
      try {
        console.log('📡 MarketWorking: Calling API...');
        
        const [listingsResult, categoriesResult] = await Promise.allSettled([
          apiClient.getListings(),
          apiClient.getCategories(),
        ]);
        
        console.log('📊 MarketWorking: Results:', { listingsResult, categoriesResult });
        
        if (listingsResult.status === 'fulfilled') {
          setListings(listingsResult.value || []);
          console.log('✅ MarketWorking: Listings loaded:', listingsResult.value?.length || 0);
        } else {
          console.error('❌ MarketWorking: Listings failed:', listingsResult.reason);
        }
        
        if (categoriesResult.status === 'fulfilled') {
          setCategories(categoriesResult.value || []);
          console.log('✅ MarketWorking: Categories loaded:', categoriesResult.value?.length || 0);
        } else {
          console.error('❌ MarketWorking: Categories failed:', categoriesResult.reason);
        }
        
        if (listingsResult.status === 'rejected' && categoriesResult.status === 'rejected') {
          setError('無法載入數據');
        }
        
      } catch (err) {
        console.error('💥 MarketWorking: Unexpected error:', err);
        setError(err instanceof Error ? err.message : '載入失敗');
      } finally {
        console.log('🏁 MarketWorking: Setting loading to false');
        setLoading(false);
      }
    };

    fetchData();
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
          <p className="text-gray-600">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="bg-white shadow-sm border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
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

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-2xl font-bold text-gray-900">商業機會列表</h2>
            <div className="text-sm text-gray-600">
              顯示 {filteredListings.length} 個結果
            </div>
          </div>

          <div className="flex flex-wrap gap-2 mb-6">
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
            {categories.map((category) => (
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

        {filteredListings.length === 0 ? (
          <div className="text-center py-12">
            <div className="text-gray-500 text-lg">暫無商業機會</div>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            {filteredListings.map((listing) => (
              <ListingCard key={listing.id} listing={listing} />
            ))}
          </div>
        )}

        <div className="mt-8 p-4 bg-blue-50 border border-blue-200 rounded">
          <p className="text-sm text-gray-600">
            <strong>✅ Pages Router 版本：</strong>這是使用 Pages Router 實現的版本，React 水合功能正常工作
          </p>
        </div>
      </div>
    </div>
  );
}
