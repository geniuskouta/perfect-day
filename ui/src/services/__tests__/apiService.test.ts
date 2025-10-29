import { describe, it, expect, vi, beforeEach } from 'vitest';
import { apiService } from '../apiService';

// Mock fetch
global.fetch = vi.fn();

describe('apiService', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Authentication', () => {
    it('should login with username', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { user: { username: 'testuser' } },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.auth.login('testuser');

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ username: 'testuser' }),
      });
      expect(result.success).toBe(true);
    });

    it('should handle login errors', async () => {
      const mockResponse = {
        ok: false,
        status: 401,
        json: vi.fn().mockResolvedValue({
          error: { code: 'UNAUTHORIZED', message: 'Invalid credentials' },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.auth.login('invaliduser');

      expect(result.success).toBe(false);
      expect(result.error).toEqual({
        code: 'UNAUTHORIZED',
        message: 'Invalid credentials',
      });
    });
  });

  describe('Perfect Days', () => {
    const mockPerfectDay = {
      title: 'Test Day',
      description: 'Test description',
      date: '2025-01-15',
      activities: [],
    };

    it('should create a perfect day', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { id: 'pd1', ...mockPerfectDay },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.create(mockPerfectDay);

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/perfect-days', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(mockPerfectDay),
      });
      expect(result.success).toBe(true);
      expect(result.data).toEqual({ id: 'pd1', ...mockPerfectDay });
    });

    it('should get a perfect day by ID', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { id: 'pd1', ...mockPerfectDay },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.getById('pd1');

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/perfect-days/pd1', {
        method: 'GET',
        credentials: 'include',
      });
      expect(result.success).toBe(true);
    });

    it('should list perfect days with pagination', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: {
            perfect_days: [{ id: 'pd1', ...mockPerfectDay }],
            pagination: { total: 1, limit: 10, offset: 0 },
          },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.list({ limit: 10, offset: 0 });

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/perfect-days?limit=10&offset=0',
        {
          method: 'GET',
          credentials: 'include',
        }
      );
      expect(result.success).toBe(true);
    });

    it('should update a perfect day', async () => {
      const updatedData = { ...mockPerfectDay, title: 'Updated Title' };
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { id: 'pd1', ...updatedData },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.update('pd1', updatedData);

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/perfect-days/pd1', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(updatedData),
      });
      expect(result.success).toBe(true);
    });

    it('should delete a perfect day', async () => {
      const mockResponse = {
        ok: true,
        status: 204,
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.delete('pd1');

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/perfect-days/pd1', {
        method: 'DELETE',
        credentials: 'include',
      });
      expect(result.success).toBe(true);
    });
  });

  describe('Places Search', () => {
    it('should search places with query', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: {
            places: [
              { name: 'Starbucks Shibuya', area: 'Shibuya' },
              { name: 'Blue Bottle Coffee', area: 'Shibuya' },
            ],
            query: 'coffee',
          },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.places.search('coffee');

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/places/search?q=coffee',
        {
          method: 'GET',
          credentials: 'include',
        }
      );
      expect(result.success).toBe(true);
      expect(result.data.places).toHaveLength(2);
    });

    it('should handle places search with limit', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { places: [], query: 'restaurant', limit: 5 },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.places.search('restaurant', { limit: 5 });

      expect(fetch).toHaveBeenCalledWith(
        'http://localhost:8080/api/v1/places/search?q=restaurant&limit=5',
        {
          method: 'GET',
          credentials: 'include',
        }
      );
    });
  });

  describe('Areas', () => {
    it('should get list of areas', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({
          success: true,
          data: { areas: ['Shibuya', 'Shinjuku', 'Harajuku'] },
        }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.areas.list();

      expect(fetch).toHaveBeenCalledWith('http://localhost:8080/api/v1/areas', {
        method: 'GET',
        credentials: 'include',
      });
      expect(result.success).toBe(true);
      expect(result.data.areas).toEqual(['Shibuya', 'Shinjuku', 'Harajuku']);
    });
  });

  describe('Error Handling', () => {
    it('should handle network errors', async () => {
      (fetch as any).mockRejectedValue(new Error('Network error'));

      const result = await apiService.perfectDays.getById('pd1');

      expect(result.success).toBe(false);
      expect(result.error.message).toBe('Network error');
    });

    it('should handle non-JSON responses', async () => {
      const mockResponse = {
        ok: false,
        status: 500,
        json: vi.fn().mockRejectedValue(new Error('Invalid JSON')),
        text: vi.fn().mockResolvedValue('Internal Server Error'),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      const result = await apiService.perfectDays.getById('pd1');

      expect(result.success).toBe(false);
      expect(result.error.message).toBe('Internal Server Error');
    });
  });

  describe('Request Configuration', () => {
    it('should include credentials in all requests', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({ success: true, data: {} }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      await apiService.perfectDays.list();

      expect(fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({ credentials: 'include' })
      );
    });

    it('should set correct content-type for POST requests', async () => {
      const mockResponse = {
        ok: true,
        json: vi.fn().mockResolvedValue({ success: true, data: {} }),
      };
      (fetch as any).mockResolvedValue(mockResponse);

      await apiService.perfectDays.create({
        title: 'Test',
        description: 'Test',
        date: '2025-01-15',
        activities: [],
      });

      expect(fetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: { 'Content-Type': 'application/json' },
        })
      );
    });
  });
});