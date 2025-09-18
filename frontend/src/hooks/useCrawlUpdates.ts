import { useEffect, useRef } from 'react';
import { useQueryClient } from '@tanstack/react-query';

interface SSENotification {
  type: string;
  url_id?: string;
  user_id?: string;
  timestamp: string;
}

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1';

export const useCrawlUpdates = (token: string | null) => {
  const queryClient = useQueryClient();
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    if (!token) return;

    // Create EventSource connection using cookies for authentication
    // withCredentials: true ensures cookies are sent with the request
    const url = `${API_BASE_URL}/crawl/stream`;
    console.log('Attempting to connect to SSE endpoint with cookies:', url);
    
    const eventSource = new EventSource(url, {
        withCredentials: true,
    });

    eventSourceRef.current = eventSource;

    eventSource.onopen = () => {
      console.log('SSE connection opened successfully');
    };

    eventSource.onmessage = (event) => {
      try {
        const notification: SSENotification = JSON.parse(event.data);
        console.log('Received SSE notification:', notification);

        switch (notification.type) {
          case 'connection':
            console.log('SSE connection confirmed');
            break;
          case 'crawl_update':
            if (notification.url_id) {
              console.log(`Invalidating queries for URL ID: ${notification.url_id}`);
              
              // Invalidate dashboard queries to refresh the URL list with updated crawl status
              queryClient.invalidateQueries({ queryKey: ['dashboard'] });
              
              // Also invalidate any other related queries if they exist
              queryClient.invalidateQueries({ 
                queryKey: ['url', notification.url_id] 
              });
              
              queryClient.invalidateQueries({ 
                queryKey: ['crawls', notification.url_id] 
              });
            }
            break;
          case 'ping':
            // Keep-alive ping, no action needed
            console.log('Received ping from server');
            break;
          default:
            console.log('Unknown SSE notification type:', notification.type);
        }
      } catch (error) {
        console.error('Error parsing SSE message:', error);
      }
    };

    eventSource.onerror = (error) => {
      console.error('SSE connection error:', error);
      console.error('EventSource readyState:', eventSource.readyState);
      console.error('EventSource url:', eventSource.url);
      
      // Attempt to reconnect after a delay
      setTimeout(() => {
        if (eventSourceRef.current?.readyState === EventSource.CLOSED) {
          console.log('Attempting to reconnect SSE...');
          // The useEffect will handle creating a new connection
        }
      }, 5000);
    };

    return () => {
      if (eventSourceRef.current) {
        console.log('Closing SSE connection');
        eventSourceRef.current.close();
        eventSourceRef.current = null;
      }
    };
  }, [token, queryClient]);

  return {
    isConnected: eventSourceRef.current?.readyState === EventSource.OPEN,
  };
};