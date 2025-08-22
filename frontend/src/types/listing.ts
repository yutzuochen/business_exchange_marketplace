export interface Listing {
  id: number;
  title: string;
  description: string;
  price: number;
  category: string;
  condition: string;
  location: string;
  status: string;
  owner_id: number;
  view_count: number;
  created_at: string;
  updated_at: string;
  brand_story?: string;
  rent?: number;
  floor?: number;
  equipment?: string;
  decoration?: string;
  annual_revenue?: number;
  gross_profit_rate?: number;
  fastest_moving_date?: string;
  phone_number?: string;
  square_meters?: number;
  industry?: string;
  deposit?: number;
  owner?: User;
  images?: Image[];
  favorites?: Favorite[];
}

export interface User {
  id: number;
  email: string;
  username: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  role: string;
  is_active: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}

export interface Image {
  id: number;
  listing_id: number;
  filename: string;
  url: string;
  alt_text?: string;
  order: number;
  is_primary: boolean;
  created_at: string;
  updated_at: string;
}

export interface Favorite {
  id: number;
  user_id: number;
  listing_id: number;
  created_at: string;
  updated_at: string;
}
