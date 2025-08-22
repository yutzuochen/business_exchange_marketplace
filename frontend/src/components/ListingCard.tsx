'use client';

import { Listing } from '@/types/listing';
import Link from 'next/link';
import Image from 'next/image';

interface ListingCardProps {
  listing: Listing;
}

export default function ListingCard({ listing }: ListingCardProps) {
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(price / 100); // 假設價格以分為單位存儲
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-TW');
  };

  return (
    <Link href={`/market/listings/${listing.id}`}>
      <div className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 overflow-hidden">
        {/* 圖片區域 */}
        <div className="aspect-[4/3] bg-gray-100 relative">
          {listing.images && listing.images.length > 0 ? (
            <Image
              src={listing.images[0].url}
              alt={listing.images[0].alt_text || listing.title}
              fill
              className="object-cover"
            />
          ) : (
            <div className="flex items-center justify-center text-gray-400">
              <span>無圖片</span>
            </div>
          )}
          {/* 狀態標籤 */}
          <div className="absolute top-2 right-2">
            <span className={`px-2 py-1 text-xs font-medium rounded-full ${
              listing.status === 'active' 
                ? 'bg-green-100 text-green-800' 
                : 'bg-gray-100 text-gray-800'
            }`}>
              {listing.status === 'active' ? '活躍' : '非活躍'}
            </span>
          </div>
        </div>

        {/* 內容區域 */}
        <div className="p-4">
          <h3 className="text-lg font-semibold text-gray-900 mb-2 line-clamp-2">
            {listing.title}
          </h3>
          
          <p className="text-gray-600 text-sm mb-3 line-clamp-2">
            {listing.description}
          </p>

          {/* 價格 */}
          <div className="text-2xl font-bold text-red-600 mb-3">
            {formatPrice(listing.price)}
          </div>

          {/* 基本信息 */}
          <div className="flex items-center justify-between text-sm text-gray-500 mb-3">
            <span>{listing.category}</span>
            <span>{listing.condition}</span>
          </div>

          {/* 位置 */}
          <div className="flex items-center text-sm text-gray-500 mb-3">
            <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            {listing.location}
          </div>

          {/* 底部信息 */}
          <div className="flex items-center justify-between text-xs text-gray-400">
            <span>瀏覽: {listing.view_count}</span>
            <span>{formatDate(listing.created_at)}</span>
          </div>
        </div>
      </div>
    </Link>
  );
}
