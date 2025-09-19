import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { dashboardApi } from '../services/dashboardApi';
import type { DashboardFilters, BatchCrawlRequest } from '../types/dashboard';

export const useDashboardData = (filters: DashboardFilters) => {
  return useQuery({
    queryKey: ['dashboard', filters],
    queryFn: () => dashboardApi.getDashboardData(filters),
    staleTime: 30000, // 30 seconds
  });
};

export const useCreateUrl = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: dashboardApi.createUrl,
    onSuccess: () => {
      // Only invalidate queries, don't automatically refetch
      // This lets the component control when to refetch
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};

export const useDeleteUrl = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: dashboardApi.deleteUrl,
    onSuccess: () => {
      // Only invalidate queries, don't automatically refetch
      // This lets the component control when to refetch
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};

export const useBatchCrawl = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: (data: BatchCrawlRequest) => dashboardApi.batchCrawl(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};

export const useStartCrawl = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: dashboardApi.startCrawl,
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};

export const useStopCrawl = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: dashboardApi.stopCrawl,
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};

export const useBatchDelete = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: dashboardApi.batchDelete,
    onSuccess: () => {
      queryClient.invalidateQueries({ 
        queryKey: ['dashboard'],
        refetchType: 'none'
      });
    },
  });
};