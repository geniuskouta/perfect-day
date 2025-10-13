const API_BASE_URL = 'http://localhost:8080/api/v1';

interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: {
    code: string;
    message: string;
  };
}

interface PerfectDay {
  id?: string;
  title: string;
  description: string;
  date: string;
  username?: string;
  activities: Activity[];
}

interface Activity {
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

interface PaginationParams {
  limit?: number;
  offset?: number;
  user?: string;
}

interface PlacesSearchParams {
  limit?: number;
}

class ApiService {
  // Transform API response from snake_case to camelCase
  private transformPerfectDayFromApi(apiData: any): PerfectDay {
    return {
      ...apiData,
      activities: apiData.activities?.map((activity: any) => ({
        ...activity,
        startTime: activity.start_time,
        duration: activity.duration_minutes || activity.duration,
      })) || []
    };
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const url = `${API_BASE_URL}${endpoint}`;
      const hasBody = options.body !== undefined;
      const headers: Record<string, string> = {};

      // Copy headers from options if they exist
      if (options.headers) {
        if (options.headers instanceof Headers) {
          options.headers.forEach((value, key) => {
            headers[key] = value;
          });
        } else if (Array.isArray(options.headers)) {
          options.headers.forEach(([key, value]) => {
            headers[key] = value;
          });
        } else {
          Object.assign(headers, options.headers);
        }
      }

      if (hasBody) {
        headers['Content-Type'] = 'application/json';
      }

      // Build request options
      const requestOptions: RequestInit = {
        credentials: 'include',
        ...options,
      };

      // Only add headers if they're not empty or if we have a body
      const hasHeaders = Object.keys(headers).length > 0;
      if (hasHeaders) {
        requestOptions.headers = headers;
      }

      // Explicitly set method to GET if no method is specified
      if (!requestOptions.method) {
        requestOptions.method = 'GET';
      }

      const response = await fetch(url, requestOptions);

      if (response.status === 204) {
        return { success: true };
      }

      let data;
      try {
        data = await response.json();
      } catch {
        const text = await response.text();
        return {
          success: false,
          error: {
            code: 'PARSE_ERROR',
            message: text || 'Failed to parse response',
          },
        };
      }

      if (!response.ok) {
        return {
          success: false,
          error: data.error || {
            code: 'HTTP_ERROR',
            message: `HTTP ${response.status}`,
          },
        };
      }

      return {
        success: true,
        data: data.data || data,
      };
    } catch (error) {
      return {
        success: false,
        error: {
          code: 'NETWORK_ERROR',
          message: error instanceof Error ? error.message : 'Network error',
        },
      };
    }
  }

  auth = {
    login: async (username: string): Promise<ApiResponse<{ user: any }>> => {
      return this.request('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ username }),
      });
    },

    logout: async (): Promise<ApiResponse> => {
      return this.request('/auth/logout', {
        method: 'POST',
      });
    },

    getCurrentUser: async (): Promise<ApiResponse<{ user: any }>> => {
      return this.request('/auth/me');
    },
  };

  perfectDays = {
    create: async (perfectDay: Omit<PerfectDay, 'id'>): Promise<ApiResponse<PerfectDay>> => {
      // Transform camelCase to snake_case for API
      const transformedPerfectDay = {
        ...perfectDay,
        activities: perfectDay.activities.map(activity => ({
          ...activity,
          start_time: activity.startTime,
          duration: activity.duration,
        }))
      };

      const result = await this.request('/perfect-days', {
        method: 'POST',
        body: JSON.stringify(transformedPerfectDay),
      });

      // Transform response from snake_case to camelCase
      if (result.success && result.data) {
        result.data = this.transformPerfectDayFromApi(result.data);
      }

      return result;
    },

    getById: async (id: string): Promise<ApiResponse<PerfectDay>> => {
      const result = await this.request(`/perfect-days/${id}`);

      // Transform response from snake_case to camelCase
      if (result.success && result.data) {
        result.data = this.transformPerfectDayFromApi(result.data);
      }

      return result;
    },

    list: async (params: PaginationParams = {}): Promise<
      ApiResponse<{
        perfect_days: PerfectDay[];
        pagination: {
          total: number;
          limit: number;
          offset: number;
        };
      }>
    > => {
      const searchParams = new URLSearchParams();
      if (params.limit) searchParams.append('limit', params.limit.toString());
      if (params.offset !== undefined) searchParams.append('offset', params.offset.toString());
      if (params.user) searchParams.append('user', params.user);

      const queryString = searchParams.toString();
      const endpoint = queryString ? `/perfect-days?${queryString}` : '/perfect-days';

      const result = await this.request(endpoint);

      // Transform response from snake_case to camelCase
      if (result.success && result.data?.perfect_days) {
        result.data.perfect_days = result.data.perfect_days.map((pd: any) =>
          this.transformPerfectDayFromApi(pd)
        );
      }

      return result;
    },

    update: async (
      id: string,
      perfectDay: Partial<PerfectDay>
    ): Promise<ApiResponse<PerfectDay>> => {
      // Transform camelCase to snake_case for API
      const transformedPerfectDay = {
        ...perfectDay,
        activities: perfectDay.activities?.map(activity => ({
          ...activity,
          start_time: activity.startTime,
          duration: activity.duration,
        })) || []
      };

      const result = await this.request<PerfectDay>(`/perfect-days/${id}`, {
        method: 'PUT',
        body: JSON.stringify(transformedPerfectDay),
      });

      // Transform response from snake_case to camelCase
      if (result.success && result.data) {
        result.data = this.transformPerfectDayFromApi(result.data);
      }

      return result;
    },

    delete: async (id: string): Promise<ApiResponse> => {
      return this.request(`/perfect-days/${id}`, {
        method: 'DELETE',
      });
    },
  };

  users = {
    getProfile: async (username: string): Promise<ApiResponse<any>> => {
      return this.request(`/users/${username}`);
    },

    getPerfectDays: async (
      username: string,
      params: PaginationParams = {}
    ): Promise<ApiResponse<any>> => {
      const searchParams = new URLSearchParams();
      if (params.limit) searchParams.append('limit', params.limit.toString());
      if (params.offset !== undefined) searchParams.append('offset', params.offset.toString());

      const queryString = searchParams.toString();
      const endpoint = queryString
        ? `/users/${username}/perfect-days?${queryString}`
        : `/users/${username}/perfect-days`;

      return this.request(endpoint);
    },
  };

  places = {
    search: async (
      query: string,
      params: PlacesSearchParams = {}
    ): Promise<
      ApiResponse<{
        places: Array<{
          name: string;
          area: string;
        }>;
        query: string;
        limit?: number;
      }>
    > => {
      const searchParams = new URLSearchParams();
      searchParams.append('q', query);
      if (params.limit) searchParams.append('limit', params.limit.toString());

      return this.request(`/places/search?${searchParams.toString()}`);
    },
  };

  areas = {
    list: async (): Promise<
      ApiResponse<{
        areas: string[];
      }>
    > => {
      return this.request('/areas');
    },
  };
}

export const apiService = new ApiService();
export type { ApiResponse, PerfectDay, Activity, PaginationParams };