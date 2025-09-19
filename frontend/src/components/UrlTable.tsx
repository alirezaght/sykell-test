import React from 'react';
import type { CrawlResult, SortColumn } from '../types/dashboard';

interface SortInfo {
  column?: SortColumn;
  order?: 'asc' | 'desc';
}

interface UrlTableProps {
  data: CrawlResult[];
  sortableColumns: Array<{ key: SortColumn; label: string }>;
  currentSort: SortInfo;
  onSort: (column: SortColumn) => void;
  onDelete: (urlId: string) => void;
  onBatchDelete: (urlIds: string[]) => void;
  onStartCrawl: (urlId: string) => void;
  onStopCrawl: (urlId: string) => void;
  onBatchCrawl: (urlIds: string[], action: 'start' | 'stop') => void;
  isDeleting: boolean;
  isBatchDeleting: boolean;
  isCrawlLoading: boolean;
  selectedUrls: string[];
  onSelectUrl: (urlId: string) => void;
  onSelectAll: (selected: boolean) => void;
}

export const UrlTable: React.FC<UrlTableProps> = ({
  data,
  sortableColumns,
  currentSort,
  onSort,
  onDelete,
  onBatchDelete,
  onStartCrawl,
  onStopCrawl,
  onBatchCrawl,
  isDeleting,
  isBatchDeleting,
  isCrawlLoading,
  selectedUrls,
  onSelectUrl,
  onSelectAll,
}) => {
  const formatDate = (dateString?: string | null) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusColor = (status?: string | null) => {
    switch (status) {
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'running':
        return 'bg-blue-100 text-blue-800';
      case 'queued':
        return 'bg-yellow-100 text-yellow-800';
      case 'failed':
        return 'bg-red-100 text-red-800';
      case 'cancelled':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  const isAllSelected = data.length > 0 && selectedUrls.length === data.length;
  const isIndeterminate = selectedUrls.length > 0 && selectedUrls.length < data.length;

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    onSelectAll(e.target.checked);
  };

  const handleSelectUrl = (urlId: string) => {
    onSelectUrl(urlId);
  };

  const canStartCrawl = (status?: string | null) => {
    return !status || status === 'completed' || status === 'failed' || status === 'cancelled';
  };

  const canStopCrawl = (status?: string | null) => {
    return status === 'running' || status === 'queued';
  };

  const SortIcon: React.FC<{ column: SortColumn }> = ({ column }) => {
    if (currentSort.column !== column) {
      return (
        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
        </svg>
      );
    }

    return currentSort.order === 'asc' ? (
      <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
      </svg>
    ) : (
      <svg className="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
      </svg>
    );
  };

  if (data.length === 0) {
    return (
      <div className="text-center py-12 bg-gray-50 rounded-lg">
        <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <h3 className="mt-2 text-sm font-medium text-gray-900">No URLs found</h3>
        <p className="mt-1 text-sm text-gray-500">Get started by adding your first URL.</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {/* Batch Controls */}
      {selectedUrls.length > 0 && (
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-blue-900">
              {selectedUrls.length} URL{selectedUrls.length > 1 ? 's' : ''} selected
            </span>
            <div className="flex space-x-2">
              <button
                onClick={() => onBatchCrawl(selectedUrls, 'start')}
                disabled={isCrawlLoading || isBatchDeleting}
                className="px-3 py-1 text-sm bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Start Crawl
              </button>
              <button
                onClick={() => onBatchCrawl(selectedUrls, 'stop')}
                disabled={isCrawlLoading || isBatchDeleting}
                className="px-3 py-1 text-sm bg-orange-600 text-white rounded hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Stop Crawl
              </button>
              <button
                onClick={() => {
                  if (confirm(`Are you sure you want to delete ${selectedUrls.length} selected URL${selectedUrls.length > 1 ? 's' : ''}?`)) {
                    onBatchDelete(selectedUrls);
                  }
                }}
                disabled={isDeleting || isBatchDeleting || isCrawlLoading}
                className="px-3 py-1 text-sm bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isBatchDeleting ? 'Deleting...' : 'Delete Selected'}
              </button>
            </div>
          </div>
        </div>
      )}

      <div className="overflow-x-auto bg-white shadow-sm rounded-lg">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left">
                <input
                  type="checkbox"
                  checked={isAllSelected}
                  ref={(input) => {
                    if (input) input.indeterminate = isIndeterminate;
                  }}
                  onChange={handleSelectAll}
                  className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
              </th>
              {sortableColumns.map(({ key, label }) => (
                <th
                  key={key}
                  onClick={() => onSort(key)}
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 select-none"
                >
                  <div className="flex items-center space-x-1">
                    <span>{label}</span>
                    <SortIcon column={key} />
                  </div>
                </th>
              ))}
              
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {data.map((item) => (
              <tr key={item.url_id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    checked={selectedUrls.includes(item.url_id)}
                    onChange={() => handleSelectUrl(item.url_id)}
                    className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm text-blue-600 hover:text-blue-800">
                    <a href={item.normalized_url} target="_blank" rel="noopener noreferrer" className="truncate max-w-xs block">
                      {item.normalized_url}
                    </a>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {item.domain}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 max-w-xs truncate">
                  {item.page_title || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`px-2 py-1 text-xs font-medium rounded-full ${getStatusColor(item.status)}`}>
                    {item.status || 'Not crawled'}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {item.html_version || '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.internal_links_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.external_links_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.inaccessible_links_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h1_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h2_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h3_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h4_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h5_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.h6_count ?? '-'}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                  {item.has_login_form === true ? (
                    <span className="text-green-600 font-medium">Yes</span>
                  ) : item.has_login_form === false ? (
                    <span className="text-gray-600">No</span>
                  ) : (
                    '-'
                  )}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {formatDate(item.url_created_at)}
                </td>                
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};