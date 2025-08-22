'use client';

import { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import { Listing } from '@/types/listing';
import { apiClient } from '@/lib/api';
import Link from 'next/link';

export default function ListingDetailPage() {
  const params = useParams();
  const id = params?.id as string;
  const [listing, setListing] = useState<Listing | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchListing = async () => {
      try {
        setLoading(true);
        setError(null);
        const data = await apiClient.getListing(parseInt(id));
        setListing(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : '載入失敗');
      } finally {
        setLoading(false);
      }
    };

    if (id) {
      fetchListing();
    }
  }, [id]);

  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(price / 100);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-TW');
  };

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
          <h1 className="text-2xl font-bold text-red-600 mb-4">載入失敗</h1>
          <p className="text-gray-600 mb-4">{error}</p>
          <Link href="/market" className="text-blue-600 hover:underline">
            返回市場
          </Link>
        </div>
      </div>
    );
  }

  if (!listing) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-gray-600 mb-4">商品不存在</h1>
          <Link href="/market" className="text-blue-600 hover:underline">
            返回市場
          </Link>
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
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{listing.title}</h1>
              <p className="mt-2 text-gray-600">商品詳情</p>
            </div>
            <div className="flex items-center space-x-4">
              <div className="flex items-center space-x-2">
                <h2 className="text-4xl font-bold text-orange-500 drop-shadow-lg">567</h2>
                <span className="text-xl font-bold text-orange-500">我來接</span>
              </div>
              <Link 
                href="/market"
                className="bg-gray-600 text-white px-4 py-2 rounded hover:bg-gray-700"
              >
                返回市場
              </Link>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* 主要內容區域 */}
          <div className="lg:col-span-2">
            {/* 商品圖片 */}
            <div className="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4 text-gray-900">商品圖片</h2>
              <div className="aspect-[4/3] bg-gray-100 rounded-lg flex items-center justify-center">
                {listing.images && listing.images.length > 0 ? (
                  <img
                    src={listing.images[0].url}
                    alt={listing.images[0].alt_text || listing.title}
                    className="w-full h-full object-cover rounded-lg"
                  />
                ) : (
                  <span className="text-gray-500">{listing.title}</span>
                )}
              </div>
            </div>

            {/* 商品描述 */}
            <div className="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4 text-gray-900">商品描述</h2>
              <p className="text-gray-900 leading-relaxed">{listing.description}</p>
            </div>

            {/* 品牌故事 */}
            {listing.brand_story && (
              <div className="bg-white rounded-lg shadow-md p-6 mb-6">
                <h2 className="text-xl font-semibold mb-4 text-gray-900">品牌故事</h2>
                <p className="text-gray-900 leading-relaxed">{listing.brand_story}</p>
              </div>
            )}
          </div>

          {/* 側邊欄 */}
          <div className="lg:col-span-1">
            {/* 價格和基本信息 */}
            <div className="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4 text-gray-900">基本信息</h2>
              <div className="space-y-3">
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">價格：</span>
                  <span className="font-semibold text-2xl text-green-600">
                    {formatPrice(listing.price)}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">分類：</span>
                  <span className="font-semibold text-gray-900">{listing.category}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">狀態：</span>
                  <span className="font-semibold text-gray-900">{listing.status}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">位置：</span>
                  <span className="font-semibold text-gray-900">{listing.location}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">瀏覽次數：</span>
                  <span className="font-semibold text-gray-900">{listing.view_count}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-gray-700 font-medium">創建時間：</span>
                  <span className="font-semibold text-gray-900">{formatDate(listing.created_at)}</span>
                </div>
              </div>
            </div>

            {/* 詳細屬性 */}
            <div className="bg-white rounded-lg shadow-md p-6 mb-6">
              <h2 className="text-xl font-semibold mb-4 text-gray-900">詳細屬性</h2>
              <div className="space-y-3">
                {listing.condition && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">狀況：</span>
                    <span className="font-semibold text-gray-900">{listing.condition}</span>
                  </div>
                )}
                {listing.industry && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">行業：</span>
                    <span className="font-semibold text-gray-900">{listing.industry}</span>
                  </div>
                )}
                {listing.square_meters && listing.square_meters > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">坪數：</span>
                    <span className="font-semibold text-gray-900">{listing.square_meters} 坪</span>
                  </div>
                )}
                {listing.rent && listing.rent > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">租金：</span>
                    <span className="font-semibold text-gray-900">{formatPrice(listing.rent)}</span>
                  </div>
                )}
                {listing.floor && listing.floor > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">樓層：</span>
                    <span className="font-semibold text-gray-900">{listing.floor} 樓</span>
                  </div>
                )}
                {listing.equipment && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">設備：</span>
                    <span className="font-semibold text-gray-900">{listing.equipment}</span>
                  </div>
                )}
                {listing.decoration && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">裝修：</span>
                    <span className="font-semibold text-gray-900">{listing.decoration}</span>
                  </div>
                )}
                {listing.annual_revenue && listing.annual_revenue > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">年營收：</span>
                    <span className="font-semibold text-gray-900">{formatPrice(listing.annual_revenue)}</span>
                  </div>
                )}
                {listing.gross_profit_rate && listing.gross_profit_rate > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">毛利率：</span>
                    <span className="font-semibold text-gray-900">{(listing.gross_profit_rate * 100).toFixed(1)}%</span>
                  </div>
                )}
                {listing.fastest_moving_date && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">最快成交日：</span>
                    <span className="font-semibold text-gray-900">{formatDate(listing.fastest_moving_date)}</span>
                  </div>
                )}
                {listing.phone_number && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">聯絡電話：</span>
                    <span className="font-semibold text-gray-900">{listing.phone_number}</span>
                  </div>
                )}
                {listing.deposit && listing.deposit > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">押金：</span>
                    <span className="font-semibold text-gray-900">{formatPrice(listing.deposit)}</span>
                  </div>
                )}
              </div>
            </div>

            {/* 賣家信息 */}
            {listing.owner && (
              <div className="bg-white rounded-lg shadow-md p-6">
                <h2 className="text-xl font-semibold mb-4 text-gray-900">賣家信息</h2>
                <div className="space-y-3">
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">用戶名：</span>
                    <span className="font-semibold text-gray-900">{listing.owner.username}</span>
                  </div>
                  {listing.owner.first_name && listing.owner.last_name && (
                    <div className="flex justify-between">
                      <span className="text-gray-700 font-medium">姓名：</span>
                      <span className="font-semibold text-gray-900">
                        {listing.owner.first_name} {listing.owner.last_name}
                      </span>
                    </div>
                  )}
                  <div className="flex justify-between">
                    <span className="text-gray-700 font-medium">註冊時間：</span>
                    <span className="font-semibold text-gray-900">{formatDate(listing.owner.created_at)}</span>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
