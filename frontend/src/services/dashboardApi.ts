import api from './api';
import type { DashboardResponse, DashboardFilters, CreateUrlRequest, CreateUrlResponse, BatchCrawlRequest, BatchCrawlResponse } from '../types/dashboard';

export const dashboardApi = {
  getDashboardData: async (filters: DashboardFilters = {}): Promise<DashboardResponse> => {
    const params = new URLSearchParams();
    
    if (filters.query_filter) {
      params.append('query', filters.query_filter);
    }
    if (filters.sort_by) {
      params.append('sort_by', filters.sort_by);
    }
    if (filters.sort_order) {
      params.append('order', filters.sort_order);
    }
    if (filters.page) {
      params.append('page', filters.page.toString());
    }
    if (filters.page_size) {
      params.append('limit', filters.page_size.toString());
    }

    const response = await api.get(`/urls?${params.toString()}`);
    
    // Transform backend response to match frontend interface
    return {
      urls: response.data.urls || [],
      total_count: response.data.total_count || 0,
      page: response.data.page || 1,
      page_size: response.data.limit || 20,
      total_pages: Math.ceil((response.data.total_count || 0) / (response.data.limit || 20)),
    };
  },

  createUrl: async (data: CreateUrlRequest): Promise<CreateUrlResponse> => {
    await api.post('/urls', { url: data.normalized_url });
    return {
      id: 'success', // Backend doesn't return ID, just success
      message: 'URL added successfully',
    };
  },

  deleteUrl: async (urlId: string): Promise<void> => {
    await api.delete(`/urls/${urlId}`);
  },

  batchCrawl: async (data: BatchCrawlRequest): Promise<BatchCrawlResponse> => {
    const { url_ids, action } = data;
    let successCount = 0;
    let failedCount = 0;
    const errors: string[] = [];

    // Process each URL individually
    for (const urlId of url_ids) {
      try {
        if (action === 'start') {
          await api.post(`/crawl/start/${urlId}`);
        } else if (action === 'stop') {
          await api.post(`/crawl/stop/${urlId}`);
        }
        successCount++;
      } catch (error) {
        failedCount++;
        errors.push(`Failed to ${action} crawl for URL ${urlId}`);
      }
    }

    return {
      success_count: successCount,
      failed_count: failedCount,
      message: `Batch ${action}: ${successCount} successful, ${failedCount} failed`,
    };
  },

  startCrawl: async (urlId: string): Promise<void> => {
    await api.post(`/crawl/start/${urlId}`);
  },

  stopCrawl: async (urlId: string): Promise<void> => {
    await api.post(`/crawl/stop/${urlId}`);
  },
};