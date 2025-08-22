import { Listing } from '@/types/listing';

// 臨時硬編碼 API URL 用於調試
const API_BASE_URL = 'https://business-exchange-backend-430730011391.us-central1.run.app';

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
      ...options,
    });

    if (!response.ok) {
      throw new Error(`API request failed: ${response.status} ${response.statusText}`);
    }

    return response.json();
  }

  // 獲取所有 listings
  async getListings(): Promise<Listing[]> {
    const response = await this.request<{listings: Listing[]}>('/api/v1/listings');
    return response.listings || [];
  }

  // 獲取單個 listing
  async getListing(id: number): Promise<Listing> {
    const response = await this.request<{listing: Listing}>('/api/v1/listings/' + id);
    return response.listing;
  }

  // 獲取 categories
  async getCategories(): Promise<string[]> {
    const response = await this.request<{categories: string[]}>('/api/v1/categories');
    return response.categories || [];
  }

  // 創建 listing
  async createListing(data: Partial<Listing>): Promise<Listing> {
    const response = await this.request<{listing: Listing}>('/api/v1/listings', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return response.listing;
  }

  // 更新 listing
  async updateListing(id: number, data: Partial<Listing>): Promise<Listing> {
    const response = await this.request<{listing: Listing}>(`/api/v1/listings/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return response.listing;
  }

  // 刪除 listing
  async deleteListing(id: number): Promise<void> {
    return this.request<void>(`/api/v1/listings/${id}`, {
      method: 'DELETE',
    });
  }
}

// 創建默認實例
export const apiClient = new ApiClient();
