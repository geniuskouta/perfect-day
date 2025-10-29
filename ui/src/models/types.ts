export interface Activity {
  id: string;
  name: string;
  location: {
    type: string;
    name: string;
    area: string;
  };
  startTime: string;
  duration: number;
  description: string;
  commentary: string;
}

export interface PerfectDay {
  id?: string;
  title: string;
  description: string;
  date: string;
  username?: string;
  activities: Activity[];
}

export interface User {
  username: string;
}

export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

export interface PaginationParams {
  limit?: number;
  offset?: number;
  user?: string;
}

export interface PlacesSearchParams {
  limit?: number;
}

export interface PlaceSearchResult {
  place_id: string;
  name: string;
  formatted_address: string;
  area?: string;
}