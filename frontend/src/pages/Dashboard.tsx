import React, { useState, useMemo } from 'react';
import { useDashboardData, useCreateUrl, useDeleteUrl, useBatchDelete, useBatchCrawl, useStartCrawl, useStopCrawl } from '../hooks/useDashboard';
import { useCrawlUpdates } from '../hooks/useCrawlUpdates';
import { useAuth } from '../context/AuthContext';
import type { DashboardFilters, SortColumn } from '../types/dashboard';
import { SearchFilter } from '../components/SearchFilter';
import { UrlTable } from '../components/UrlTable';
import { Pagination } from '../components/Pagination';
import { AddUrlModal } from '../components/AddUrlModal';
import { LoadingSpinner } from '../components/LoadingSpinner';

const DEFAULT_PAGE_SIZE = 5;

export const Dashboard: React.FC = () => {
  const { token } = useAuth();
  const [filters, setFilters] = useState<DashboardFilters>({
    page: 1,
    page_size: DEFAULT_PAGE_SIZE,
    sort_by: 'created_at',
    sort_order: 'desc',
  });
  const [isAddUrlModalOpen, setIsAddUrlModalOpen] = useState(false);
  const [selectedUrls, setSelectedUrls] = useState<string[]>([]);

  // Initialize SSE connection for real-time crawl updates
  const { isConnected } = useCrawlUpdates(token);

  const { data, isLoading, error, refetch } = useDashboardData(filters);
  const createUrlMutation = useCreateUrl();
  const deleteUrlMutation = useDeleteUrl();
  const batchDeleteMutation = useBatchDelete();
  const batchCrawlMutation = useBatchCrawl();
  const startCrawlMutation = useStartCrawl();
  const stopCrawlMutation = useStopCrawl();

  const isCrawlLoading = batchCrawlMutation.isPending || startCrawlMutation.isPending || stopCrawlMutation.isPending;

  const handleSearch = (query: string) => {
    setFilters(prev => ({
      ...prev,
      query_filter: query || undefined,
    }));
  };

  const handleSort = (column: SortColumn) => {
    setFilters(prev => ({
      ...prev,
      sort_by: column,
      sort_order: prev.sort_by === column && prev.sort_order === 'asc' ? 'desc' : 'asc',
    }));
  };

  const handlePageChange = (page: number) => {
    console.log('Changing to page:', page);
    setFilters(prev => {
      const newFilters = { ...prev, page };
      console.log('New filters:', newFilters);
      return newFilters;
    });
  };

  const handleCreateUrl = async (url: string) => {
    try {
      await createUrlMutation.mutateAsync({ normalized_url: url });
      setIsAddUrlModalOpen(false);
      // Manually refetch to ensure fresh data
      refetch();
    } catch (error) {
      console.error('Failed to create URL:', error);
      // You might want to show a toast notification here
    }
  };

  const handleDeleteUrl = async (urlId: string) => {
    if (!confirm('Are you sure you want to delete this URL?')) {
      return;
    }
    
    try {
      await deleteUrlMutation.mutateAsync(urlId);            
      // Manually refetch to ensure fresh data
      refetch();
    } catch (error) {
      console.error('Failed to delete URL:', error);
      // You might want to show a toast notification here
    }
  };

  const handleStartCrawl = async (urlId: string) => {
    try {
      await startCrawlMutation.mutateAsync(urlId);
      refetch();
    } catch (error) {
      console.error('Failed to start crawl:', error);
    }
  };

  const handleStopCrawl = async (urlId: string) => {
    try {
      await stopCrawlMutation.mutateAsync(urlId);
      refetch();
    } catch (error) {
      console.error('Failed to stop crawl:', error);
    }
  };

  const handleBatchCrawl = async (urlIds: string[], action: 'start' | 'stop') => {
    try {
      await batchCrawlMutation.mutateAsync({ url_ids: urlIds, action });
      setSelectedUrls([]); // Clear selection after batch operation
      refetch();
    } catch (error) {
      console.error(`Failed to ${action} batch crawl:`, error);
    }
  };

  const handleBatchDelete = async (urlIds: string[]) => {
    try {
      const result = await batchDeleteMutation.mutateAsync(urlIds);
      setSelectedUrls([]); // Clear selection after batch operation
      refetch();
      
      // You might want to show a toast notification here
      if (result.failedCount > 0) {
        console.warn(`Batch delete completed with errors: ${result.successCount} successful, ${result.failedCount} failed`);
        // Optionally show error details: result.errors
      } else {
        console.log(`Successfully deleted ${result.successCount} URLs`);
      }
    } catch (error) {
      console.error('Failed to batch delete URLs:', error);
    }
  };

  const handleSelectUrl = (urlId: string) => {
    setSelectedUrls(prev => 
      prev.includes(urlId) 
        ? prev.filter(id => id !== urlId)
        : [...prev, urlId]
    );
  };

  const handleSelectAll = (selected: boolean) => {
    if (selected) {
      setSelectedUrls(data?.urls.map(url => url.url_id) || []);
    } else {
      setSelectedUrls([]);
    }
  };

  const sortableColumns = useMemo(() => [
    { key: 'url' as SortColumn, label: 'URL' },
    { key: 'domain' as SortColumn, label: 'Domain' },
    { key: 'title' as SortColumn, label: 'Page Title' },
    { key: 'status' as SortColumn, label: 'Status' },
    { key: 'html_version' as SortColumn, label: 'HTML Version' },
    { key: 'internal_links' as SortColumn, label: 'Internal Links' },
    { key: 'external_links' as SortColumn, label: 'External Links' },
    { key: 'inaccessible_links' as SortColumn, label: 'Inaccessible Links' },
    { key: 'h1_count' as SortColumn, label: 'H1' },
    { key: 'h2_count' as SortColumn, label: 'H2' },
    { key: 'h3_count' as SortColumn, label: 'H3' },
    { key: 'h4_count' as SortColumn, label: 'H4' },
    { key: 'h5_count' as SortColumn, label: 'H5' },
    { key: 'h6_count' as SortColumn, label: 'H6' },
    { key: 'has_login_form' as SortColumn, label: 'Login Form' },
    { key: 'created_at' as SortColumn, label: 'Created At' },
  ], []);

  if (error) {
    return (
      <div className="p-4">
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <h3 className="text-lg font-medium text-red-800">Error loading dashboard</h3>
          <p className="text-red-600 mt-1">
            {error instanceof Error ? error.message : 'An unexpected error occurred'}
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-6">
        <div className="flex justify-between items-center mb-4">
          <div className="flex items-center space-x-4">
            <h1 className="text-3xl font-bold text-gray-900">URL Dashboard</h1>
            {/* SSE Connection Status Indicator */}
            <div className="flex items-center space-x-2">
              <div 
                className={`w-3 h-3 rounded-full ${
                  isConnected ? 'bg-green-500' : 'bg-red-500'
                }`}
                title={isConnected ? 'Live updates connected' : 'Live updates disconnected'}
              />
              <span className="text-sm text-gray-600">
                {isConnected ? 'Live' : 'Offline'}
              </span>
            </div>
          </div>
          <button
            onClick={() => setIsAddUrlModalOpen(true)}
            className="bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Add URL
          </button>
        </div>
        
        <SearchFilter 
          onSearch={handleSearch} 
          initialValue={filters.query_filter || ''}
        />
      </div>

      {isLoading ? (
        <LoadingSpinner />
      ) : (
        <>
          <UrlTable
            data={data?.urls || []}
            sortableColumns={sortableColumns}
            currentSort={{ column: filters.sort_by, order: filters.sort_order }}
            onSort={handleSort}
            onDelete={handleDeleteUrl}
            onBatchDelete={handleBatchDelete}
            onStartCrawl={handleStartCrawl}
            onStopCrawl={handleStopCrawl}
            onBatchCrawl={handleBatchCrawl}
            isDeleting={deleteUrlMutation.isPending}
            isBatchDeleting={batchDeleteMutation.isPending}
            isCrawlLoading={isCrawlLoading}
            selectedUrls={selectedUrls}
            onSelectUrl={handleSelectUrl}
            onSelectAll={handleSelectAll}
          />

          {data && (
            <div className="mt-6">
              <Pagination
                currentPage={filters.page || 1}
                totalPages={data.total_pages}
                totalCount={data.total_count}
                pageSize={data.page_size}
                onPageChange={handlePageChange}
              />
            </div>
          )}
        </>
      )}

      <AddUrlModal
        isOpen={isAddUrlModalOpen}
        onClose={() => setIsAddUrlModalOpen(false)}
        onSubmit={handleCreateUrl}
        isLoading={createUrlMutation.isPending}
      />
    </div>
  );
};