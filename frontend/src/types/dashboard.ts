export interface CrawlResult {
  url_id: string;
  normalized_url: string;
  domain: string;
  url_created_at: string | null;
  crawl_id?: string | null;
  status?: 'queued' | 'running' | 'completed' | 'failed' | 'cancelled' | null;
  queued_at?: string | null;
  started_at?: string | null;
  finished_at?: string | null;
  html_version?: string | null;
  page_title?: string | null;
  h1_count?: number | null;
  h2_count?: number | null;
  h3_count?: number | null;
  h4_count?: number | null;
  h5_count?: number | null;
  h6_count?: number | null;
  internal_links_count?: number | null;
  external_links_count?: number | null;
  inaccessible_links_count?: number | null;
  has_login_form?: boolean | null;
  error_message?: string | null;
  crawl_created_at?: string | null;
  crawl_updated_at?: string | null;
}

export interface DashboardResponse {
  urls: CrawlResult[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface DashboardFilters {
  query_filter?: string;
  sort_by?: SortColumn;
  sort_order?: 'asc' | 'desc';
  page?: number;
  page_size?: number;
}

export type SortColumn = 
  | 'url'
  | 'domain'
  | 'title'
  | 'status'
  | 'html_version'
  | 'internal_links'
  | 'external_links'
  | 'inaccessible_links'
  | 'h1_count'
  | 'h2_count'
  | 'h3_count'
  | 'created_at'
  | 'finished_at';

export interface CreateUrlRequest {
  normalized_url: string;
}

export interface CreateUrlResponse {
  id: string;
  message: string;
}

export interface BatchCrawlRequest {
  url_ids: string[];
  action: 'start' | 'stop';
}

export interface BatchCrawlResponse {
  success_count: number;
  failed_count: number;
  message: string;
}