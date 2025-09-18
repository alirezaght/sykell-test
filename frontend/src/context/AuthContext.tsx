import { createContext, useContext, useEffect, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { User, LoginResponse } from '../types/auth';
import type { LoginFormData } from '../types/validation';
import { authAPI } from '../services/auth';

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (credentials: LoginFormData) => Promise<LoginResponse>;
  logout: () => Promise<void>;
  token: string | null; // Keep for compatibility but will be null when using cookies
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: React.ReactNode;
}

export const AuthProvider = ({ children }: AuthProviderProps) => {
  // For cookie-based auth, we don't store tokens in localStorage
  // Instead, we'll determine auth state by attempting to fetch user profile
  const [isInitialized, setIsInitialized] = useState(false);
  const [isLoggedOut, setIsLoggedOut] = useState(false);
  const queryClient = useQueryClient();

  // Get current user query - this will work if cookie is valid
  const {
    data: user,
    isLoading,
    error,
    isSuccess,
  } = useQuery({
    queryKey: ['auth', 'user'],
    queryFn: authAPI.getCurrentUser,
    enabled: !isLoggedOut, // Don't run if user explicitly logged out
    retry: false,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });

  // Login mutation
  const loginMutation = useMutation({
    mutationFn: authAPI.login,
    onSuccess: (data) => {
      // Cookie is set by the server, just update the query cache
      setIsLoggedOut(false); // Reset logout flag
      queryClient.setQueryData(['auth', 'user'], data.user);
      // Keep token in localStorage for compatibility but primarily use cookies
      localStorage.setItem('token', data.token);
      localStorage.setItem('user', JSON.stringify(data.user));
    },
  });

  // Logout function
  const logout = async () => {
    try {
      // Call logout API to clear server-side cookie
      await authAPI.logout();
    } catch (error) {
      console.error('Logout API call failed:', error);
      // Continue with client-side cleanup even if API call fails
    } finally {
      // Clean up client-side storage
      setIsLoggedOut(true); // Set logout flag to prevent further queries
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      queryClient.clear();
    }
  };

  // Handle auth errors and initialization
  useEffect(() => {
    if (error) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      // Don't set isLoggedOut here - let user explicitly logout
    }
    if (isSuccess || error) {
      setIsInitialized(true);
    }
  }, [error, isSuccess]);

  const isAuthenticated = !!user && !isLoggedOut;

  const value: AuthContextType = {
    user: user || null,
    isLoading: isLoading && !isInitialized,
    isAuthenticated,
    login: loginMutation.mutateAsync,
    logout,
    token: localStorage.getItem('token'), // Keep for compatibility
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};